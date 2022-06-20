package http_log_server

import (
	"context"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/pawel33317/coreCommunicationFramework/db_handler"
	"github.com/pawel33317/coreCommunicationFramework/logger"
)

type HttpLogServer struct {
	closeServerChan <-chan bool
	log             logger.LogWrapper
	srv             *http.Server
	logReader       db_handler.DbLogReader
}

func NewHttpLogServer(killHttpOutChannel <-chan bool, logg logger.Logger, logsReader db_handler.DbLogReader) *HttpLogServer {
	return &HttpLogServer{closeServerChan: killHttpOutChannel, log: *logger.NewLogWrapper(logg, "HLS"), srv: nil, logReader: logsReader}
}

type Page struct {
	Title string
	Body  []byte
}

func (h *HttpLogServer) mainHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("http_log_server/logserv.html"))

	type LogPageData struct {
		PageTitle string
		LogsData  []db_handler.LogDataFormat
	}

	notParsedLogs := h.logReader.GetLogs()
	for k, elem := range notParsedLogs {
		levelVal, err := strconv.Atoi(elem.LogLevel)
		if err == nil {
			notParsedLogs[k].LogLevel = logger.LogLevel(levelVal).ToString()
		}

		timeInt, err := strconv.ParseInt(elem.LogTime, 10, 64)
		if err == nil {
			unixTimeUTC := time.Unix(timeInt, 0)
			notParsedLogs[k].LogTime = unixTimeUTC.Format("2006-01-02 15:04:05")
		}

	}

	data := LogPageData{
		PageTitle: "Log server page",
		LogsData:  notParsedLogs,
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
	}()

}
