package main

import (
	"github.com/pawel33317/coreCommunicationFramework/app_state_manager"
	"github.com/pawel33317/coreCommunicationFramework/app_state_manager/app_state"
	"github.com/pawel33317/coreCommunicationFramework/logger"
)

type AppStateClient struct {
	asm    app_state_manager.AppStateClientHandler
	Name   string
	logger logger.Logger
}

func (client *AppStateClient) OnAppStateChanged(startState app_state.State) {
	client.logger.Log(logger.DEBUG, "MAIN", "Module informed about new state", startState.ToString())
}

func (client *AppStateClient) Start(state app_state.State) {
	client.logger.Log(logger.INFO, "MAIN", "Client register as obserber and LOCK state", state.ToString())
	client.asm.RegisterObserver(client)
	client.asm.RegisterLockState(client, state)
}

func (client *AppStateClient) End(state app_state.State) {
	client.logger.Log(logger.INFO, "MAIN", "Client unlock state", state.ToString())
	client.asm.UnlockState(client)
}

func main() {
	//dbHandler.RunHandler("test")

	log := logger.NewLoggerImp()
	log.Enable()
	log.SetMinLogLevel(logger.INFO)
	asManager := app_state_manager.MakeAppStateManagerImp(log)
	asClient := &AppStateClient{asManager, "A", log}
	asClient.Start(app_state.LOADING)

	asClient2 := &AppStateClient{asManager, "B", log}
	asClient2.Start(app_state.CONFIGURED)

	asManager.Start(app_state.INITIALIZING)

	asClient.End(app_state.LOADING)
	asClient2.End(app_state.CONFIGURED)

}
