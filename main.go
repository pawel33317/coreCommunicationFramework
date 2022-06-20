package main

import (
	"github.com/pawel33317/coreCommunicationFramework/app_state_manager"
	"github.com/pawel33317/coreCommunicationFramework/db_handler"
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
	logg.SetMinLogLevel(logger.INFO)
	log := logger.NewLogWrapper(logg, "MAIN")
	log.Log(logger.INFO, "App start")
	defer log.Log(logger.INFO, "App end")

	//SM
	asManager := app_state_manager.NewAppStateManagerImp(logg)

	//TERMINATION HANDLING
	termSignalCh := make(chan bool, 1)
	sys_signal_receiver.ReceiveTerminationSignal(termSignalCh)

	//TEST CODE
	inc := 0
	if inc == 0 {
		test_objects.TestASMClient(asManager, logg)
		inc++
	}
	test_objects.Print_data_time_parallely(logg)

	//MAIN LOOP
	select {
	case termSig := <-termSignalCh:
		log.Log(logger.INFO, "Term signal received, exiting", termSig)
		//TODO: SM to shutdown
		break
	}

}
