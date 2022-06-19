package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/pawel33317/coreCommunicationFramework/app_state_manager"
	"github.com/pawel33317/coreCommunicationFramework/app_state_manager/app_state"
	"github.com/pawel33317/coreCommunicationFramework/db_handler"
	"github.com/pawel33317/coreCommunicationFramework/logger"
)

type AppStateClient struct {
	asm    app_state_manager.AppStateClientHandler
	Name   string
	logger *logger.LogWrapper
}

func (client *AppStateClient) OnAppStateChanged(startState app_state.State) {
	client.logger.Log(logger.DEBUG, "Module informed about new state", startState.ToString())
}

func (client *AppStateClient) Start(state app_state.State) {
	client.logger.Log(logger.INFO, "Client register as obserber and LOCK state", state.ToString())
	client.asm.RegisterObserver(client)
	client.asm.RegisterLockState(client, state)
}

func (client *AppStateClient) End(state app_state.State) {
	client.logger.Log(logger.INFO, "Client unlock state", state.ToString())
	client.asm.UnlockState(client)
}

func testASMClient(asManager *app_state_manager.AppStateManagerImp, log logger.Logger) {
	asClient := &AppStateClient{asManager, "A", logger.NewLogWrapper(log, "ASC1")}
	asClient.Start(app_state.LOADING)

	asClient2 := &AppStateClient{asManager, "B", logger.NewLogWrapper(log, "ASC2")}
	asClient2.Start(app_state.CONFIGURED)

	asManager.Start(app_state.INITIALIZING)

	asClient.End(app_state.LOADING)
	asClient2.End(app_state.CONFIGURED)
	log.Log(logger.INFO, "MAIN", "App end")
}

func main() {
	db := db_handler.SQLiteDb{}
	dbErr := db.Open()

	if dbErr != nil {
		panic(dbErr)
	}
	defer db.Close()

	log := logger.NewLoggerImp(&db)
	log.Enable()
	log.SetMinLogLevel(logger.DEBUG)

	log.Log(logger.INFO, "MAIN", "App start")

	asManager := app_state_manager.NewAppStateManagerImp(log)

	//napisać wątek, który sobie odlicza sekundy i drugi z odbiorem sygnału kill jak poniżej
	//czekać na info z obu i zabić apkę na ctrl+c
	inc := 0
	for {
		if inc == 0 {
			testASMClient(asManager, log)
			inc++
		}
		sigs := make(chan os.Signal, 1)
		done := make(chan bool, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			for {
				sig := <-sigs
				fmt.Println()
				fmt.Println(sig)
				done <- true
			}
		}()
		fmt.Println("awaiting signal")
		<-done
		fmt.Println("exiting")
		os.Exit(1)

	}

}
