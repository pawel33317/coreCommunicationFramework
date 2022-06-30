package main

import (
	"time"

	"github.com/pawel33317/coreCommunicationFramework/app_state_manager"
	"github.com/pawel33317/coreCommunicationFramework/app_state_manager/app_state"
	"github.com/pawel33317/coreCommunicationFramework/db_handler"
	"github.com/pawel33317/coreCommunicationFramework/dbus_server"
	"github.com/pawel33317/coreCommunicationFramework/http_log_server"
	"github.com/pawel33317/coreCommunicationFramework/logger"
	"github.com/pawel33317/coreCommunicationFramework/sys_signal_listener"
	"github.com/pawel33317/coreCommunicationFramework/tcp_server"
)

// func main() { mainthread.Init(defferedMain) } //allows to calls methods on main thread
// considere in the future

func main() {

	//DB
	db := &db_handler.SQLiteDb{}
	dbErr := db.Open()
	if dbErr != nil {
		panic(dbErr)
	}
	defer db.Close()

	//LOGGER
	loggerImp := logger.NewLoggerImp(db)
	loggerImp.Enable()
	loggerImp.SetMinLogLevel(logger.DEBUG)
	log := logger.NewLogWrapper(loggerImp, "MAIN")
	log.Log(logger.INFO, "App start")
	defer log.Log(logger.INFO, "App end")

	appStateManager := app_state_manager.NewAppStateManagerImp(loggerImp)

	termSignalChan := make(chan struct{})
	sys_signal_listener.ListenOnTerminationSignal(termSignalChan)

	http_log_server.NewHttpLogServer(loggerImp, db, appStateManager)

	tcpDataChan := make(chan string)
	tcp_server.NewTcpServer(loggerImp, appStateManager, tcpDataChan)

	dbus_server.NewDBUSServer(loggerImp, appStateManager)

	appStateManager.Start(app_state.INITIALIZED)

	timer1 := time.NewTimer(12 * time.Second)
	ticker1 := time.NewTicker(15 * time.Second)

	for {
		select {
		case termSignal := <-termSignalChan:
			log.Log(logger.WARN, "Term signal received, exiting", termSignal)
			appStateManager.RequestStateChange(app_state.SHUTDOWN)
			return
		case tcpData := <-tcpDataChan:
			log.Log(logger.INFO, "Main thread received TCP data:", tcpData)
		case <-timer1.C:
			log.Log(logger.INFO, "Timer executed")
		case <-ticker1.C:
			log.Log(logger.INFO, "Ticker executed")
		}
	}
}
