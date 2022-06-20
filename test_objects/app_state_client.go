package test_objects

import (
	"github.com/pawel33317/coreCommunicationFramework/app_state_manager"
	"github.com/pawel33317/coreCommunicationFramework/app_state_manager/app_state"
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

func TestASMClient(asManager *app_state_manager.AppStateManagerImp, log logger.Logger) {
	asClient := &AppStateClient{asManager, "A", logger.NewLogWrapper(log, "ASC1")}
	asClient.Start(app_state.LOADING)

	asClient2 := &AppStateClient{asManager, "B", logger.NewLogWrapper(log, "ASC2")}
	asClient2.Start(app_state.CONFIGURED)

	asManager.Start(app_state.INITIALIZING)

	asClient.End(app_state.LOADING)
	asClient2.End(app_state.CONFIGURED)
}
