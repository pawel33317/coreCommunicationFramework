package http_log_server

import (
	"context"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/pawel33317/coreCommunicationFramework/app_state_manager"
	"github.com/pawel33317/coreCommunicationFramework/app_state_manager/app_state"
	"github.com/pawel33317/coreCommunicationFramework/db_handler"
	"github.com/pawel33317/coreCommunicationFramework/logger"
)

type HttpLogServer struct {
	closeServerChan chan bool
	closedServerChan chan bool
	log             logger.LogWrapper
	srv             *http.Server
	logReader       db_handler.DbLogReader
	asm             app_state_manager.AppStateClientHandler
}

func NewHttpLogServer(logg logger.Logger, logsReader db_handler.DbLogReader, asmClient app_state_manager.AppStateClientHandler) *HttpLogServer {
	hls := HttpLogServer{closeServerChan: make(chan bool, 1), closedServerChan: make(chan bool, 1), log: *logger.NewLogWrapper(logg, "HLS"), srv: nil, logReader: logsReader, asm: asmClient}
	asmClient.RegisterObserver(&hls)
	return &hls
}

type Page struct {
	Title string
	Body  []byte
}

func (h *HttpLogServer) OnAppStateChanged(state app_state.State) {
	switch state {
	case app_state.INITIALIZED:
		h.asm.RegisterLockState(h, app_state.SHUTDOWN)
	case app_state.CONFIGURED:
		h.log.Log(logger.INFO, "CONFIGURED received")
		h.RunLogServer()
	case app_state.SHUTDOWN:
		h.log.Log(logger.INFO, "SHUTDOWN received")
		h.closeServerChan <- true
		<-h.closedServerChan
		h.asm.UnlockState(h)
	}
}

func (h *HttpLogServer) mainHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("http_log_server/logserv.html"))

	type LogPageData struct {
		PageTitle string
		LogsData  []db_handler.LogDataFormat
	}

	dbLogs := h.logReader.GetLogs()
	for k, elem := range dbLogs {
		levelVal, err := strconv.Atoi(elem.LogLevel)
		if err == nil {
			dbLogs[k].LogLevel = logger.LogLevel(levelVal).ToString()
		}

		timeInt, err := strconv.ParseInt(elem.LogTime, 10, 64)
		if err == nil {
			unixTimeUTC := time.Unix(timeInt, 0)
			dbLogs[k].LogTime = unixTimeUTC.Format("2006-01-02 15:04:05")
		}

	}

	data := LogPageData{
		PageTitle: "Log server page",
		LogsData:  dbLogs,
	}

	tmpl.Execute(w, data)
}

func (hls *HttpLogServer) RunLogServer() {
	hls.srv = &http.Server{Addr: ":8080"}
	http.HandleFunc("/", hls.mainHandler)

	go func() {
		go func() {
			hls.log.Log(logger.INFO, "Starting hls server")
			if err := hls.srv.ListenAndServe(); err != http.ErrServerClosed {
				hls.log.Log(logger.INFO, "Hls server error, ListenAndServe():", err)
				os.Exit(2)
			}
		}()
		<-hls.closeServerChan
		hls.log.Log(logger.INFO, "Stopping hls server")
		hls.srv.Shutdown(context.TODO())
		hls.closedServerChan<-true
	}()

}
