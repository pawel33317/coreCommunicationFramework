package http_log_server

import (
	"context"
	"encoding/json"
	"fmt"
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
	closeSignalChan chan struct{}
	log             logger.LogWrapper
	srvHttp         *http.Server
	srvHttps        *http.Server
	logReader       db_handler.DbLogReader
	asm             app_state_manager.AppStateClientHandler
}

func NewHttpLogServer(logg logger.Logger, logsReader db_handler.DbLogReader, asmClient app_state_manager.AppStateClientHandler) *HttpLogServer {
	hls := HttpLogServer{
		closeSignalChan: make(chan struct{}),
		log:             *logger.NewLogWrapper(logg, "HLS"),
		srvHttp:         nil,
		srvHttps:        nil,
		logReader:       logsReader,
		asm:             asmClient,
	}
	asmClient.RegisterObserver(&hls)
	return &hls
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
		h.closeSignalChan <- struct{}{}
		<-h.closeSignalChan
		h.asm.UnlockState(h)
	case app_state.DISABLED:
		h.log.Log(logger.INFO, "DISABLED received")
	}
}

func parseLogs(logs []db_handler.Log) {
	for k, elem := range logs {
		levelVal, err := strconv.Atoi(elem.Level)
		if err == nil {
			logs[k].Level = logger.LogLevel(levelVal).ToString()
		}

		timeInt, err := strconv.ParseInt(elem.Time, 10, 64)
		if err == nil {
			unixTimeUTC := time.Unix(timeInt, 0)
			logs[k].Time = unixTimeUTC.Format("2006-01-02 15:04:05")
		}
	}
}

func (h *HttpLogServer) mainHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("http_log_server/logserv.html"))

	type LogPageData struct {
		PageTitle string
		LogsData  []db_handler.Log
	}

	logs := h.logReader.GetLogs()
	parseLogs(logs)

	data := LogPageData{
		PageTitle: "Log server page",
		LogsData:  logs,
	}

	tmpl.Execute(w, data)
}

func logsToJson(logs []db_handler.Log) string {
	parseLogs(logs)
	jsonLogs, err := json.Marshal(logs)

	if err != nil {
		panic("Cannot convert logs to json")
	}
	return string(jsonLogs)
}

func (h *HttpLogServer) getNewLogsHandler(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["lastLogId"]
	if !ok || len(keys[0]) < 1 {
		h.log.Log(logger.WARN, "Incorrect getNewLogs request, Url param 'lastLogId' is missing", ok)
		return
	}

	lastLogId, err := strconv.Atoi(keys[0])
	if err != nil {
		h.log.Log(logger.WARN, "Incorrect getNewLogs request, Url param 'lastLogId' is not integer val", ok)
	}

	newLogs := h.logReader.GetLogsNewerThan(lastLogId)

	h.log.Log(logger.DEBUG, "Sending new logs: ", len(newLogs))
	fmt.Fprint(w, logsToJson(newLogs))
}

func (hls *HttpLogServer) RunLogServer() {
	hls.srvHttp = &http.Server{Addr: ":2001"}
	hls.srvHttps = &http.Server{Addr: ":2002"}
	http.HandleFunc("/", hls.mainHandler)
	http.HandleFunc("/getNewLogs", hls.getNewLogsHandler)

	go func() {
		go func() {
			hls.log.Log(logger.INFO, "Starting hls http server")
			if err := hls.srvHttp.ListenAndServe(); err != http.ErrServerClosed {
				hls.log.Log(logger.INFO, "Hls http server error, ListenAndServe():", err)
				os.Exit(2)
			}
		}()
		go func() {
			hls.log.Log(logger.INFO, "Starting hls https server")
			if err := hls.srvHttps.ListenAndServeTLS("test_objects/keys/server.crt", "test_objects/keys/server.key"); err != http.ErrServerClosed {
				hls.log.Log(logger.INFO, "Hls https server error, ListenAndServe():", err)
				os.Exit(2)
			}
		}()
		<-hls.closeSignalChan
		hls.log.Log(logger.INFO, "Stopping hls server")
		hls.srvHttp.Shutdown(context.TODO())
		hls.srvHttps.Shutdown(context.TODO())
		hls.closeSignalChan <- struct{}{}
	}()
}
