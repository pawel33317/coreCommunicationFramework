package main

import (
	"github.com/pawel33317/coreCommunicationFramework/app_state_manager"
	"github.com/pawel33317/coreCommunicationFramework/app_state_manager/app_state"
	"github.com/pawel33317/coreCommunicationFramework/db_handler"
	"github.com/pawel33317/coreCommunicationFramework/http_log_server"
	"github.com/pawel33317/coreCommunicationFramework/logger"
	"github.com/pawel33317/coreCommunicationFramework/sys_signal_receiver"
	"github.com/pawel33317/coreCommunicationFramework/test_objects"
)

func main() {

	//DB
	db := db_handler.SQLiteDb{}
	dbErr := db.Open()
	if dbErr != nil {
		panic(dbErr)
	}
	defer db.Close()

	//LOGGER
	logg := logger.NewLoggerImp(&db)
	logg.Enable()
	logg.SetMinLogLevel(logger.DEBUG)
	log := logger.NewLogWrapper(logg, "MAIN")
	log.Log(logger.INFO, "App start")
	defer log.Log(logger.INFO, "App end")

	//SM
	asManager := app_state_manager.NewAppStateManagerImp(logg)

	//TERMINATION HANDLING
	termSignalCh := make(chan bool, 1)
	sys_signal_receiver.ReceiveTerminationSignal(termSignalCh)

	//TEST CODE ####################<!--
	//Create ASM clients
	asClient := &test_objects.AppStateClient{Asm: asManager, Name: "A", Logger: logger.NewLogWrapper(logg, "ASC1")}
	asClient.StartClientAndLockState(app_state.LOADED)
	asClient2 := &test_objects.AppStateClient{Asm: asManager, Name: "B", Logger: logger.NewLogWrapper(logg, "ASC2")}
	asClient2.StartClientAndLockState(app_state.CONFIGURED)
	//##############################-->

	//START SM
	asManager.Start(app_state.INITIALIZED)

	//TEST CODE ####################<!--
	//Unlock states by ASM clients, start thread printing time
	asClient.UnlockState(app_state.LOADED)
	asClient2.UnlockState(app_state.CONFIGURED)
	test_objects.Print_data_time_parallely(logg)
	//##############################-->

	//HTTPS LOG SERVER
	closeServerChan := make(chan bool, 1)
	hls := http_log_server.NewHttpLogServer(closeServerChan, logg, &db)
	hls.RunLogServer()

	//MAIN LOOP
	select {
	case termSig := <-termSignalCh:
		log.Log(logger.INFO, "Term signal received, exiting", termSig)
		asManager.RequestStateChange(app_state.SHUTDOWN) //send to SM clients info about shutdown
		closeServerChan <- true                          //close https server
		break
	}

	//TODO: Make https server as SM client
}
