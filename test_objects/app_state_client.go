package test_objects

import (
	"github.com/pawel33317/coreCommunicationFramework/app_state_manager"
	"github.com/pawel33317/coreCommunicationFramework/app_state_manager/app_state"
	"github.com/pawel33317/coreCommunicationFramework/logger"
)

type AppStateClient struct {
	Asm    app_state_manager.AppStateClientHandler
	Name   string
	Logger *logger.LogWrapper
}

func (client *AppStateClient) OnAppStateChanged(startState app_state.State) {
	client.Logger.Log(logger.DEBUG, "Module informed about new state", startState.ToString())
}

func (client *AppStateClient) StartClientAndLockState(state app_state.State) {
	client.Logger.Log(logger.INFO, "Client register as obserber and LOCK state", state.ToString())
	client.Asm.RegisterObserver(client)
	client.Asm.RegisterLockState(client, state)
}

func (client *AppStateClient) UnlockState(state app_state.State) {
	client.Logger.Log(logger.INFO, "Client unlock state", state.ToString())
	client.Asm.UnlockState(client)
}
