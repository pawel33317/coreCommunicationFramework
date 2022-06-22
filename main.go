package main

import (
	"github.com/pawel33317/coreCommunicationFramework/app_state_manager"
	"github.com/pawel33317/coreCommunicationFramework/app_state_manager/app_state"
	"github.com/pawel33317/coreCommunicationFramework/db_handler"
	"github.com/pawel33317/coreCommunicationFramework/http_log_server"
	"github.com/pawel33317/coreCommunicationFramework/logger"
	"github.com/pawel33317/coreCommunicationFramework/sys_signal_listener"
	"github.com/pawel33317/coreCommunicationFramework/test_objects"
)

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

	termSignalChan := make(chan bool)
	sys_signal_listener.ListenOnTerminationSignal(termSignalChan)

	http_log_server.NewHttpLogServer(loggerImp, db, appStateManager)

	test_objects.Print_data_time_parallely(loggerImp)

	appStateManager.Start(app_state.INITIALIZED)

	select {
	case termSignal := <-termSignalChan:
		log.Log(logger.WARN, "Term signal received, exiting", termSignal)
		appStateManager.RequestStateChange(app_state.SHUTDOWN)
		break
	}
}
