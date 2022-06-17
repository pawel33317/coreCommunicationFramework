package main

import (
	"fmt"

	"github.com/pawel33317/coreCommunicationFramework/app_state_manager"
	"github.com/pawel33317/coreCommunicationFramework/app_state_manager/app_state"
)

type AppStateClient struct {
	asm  *app_state_manager.AppStateManagerCtx
	Name string
}

func (moduleData *AppStateClient) OnAppStateChanged(startState app_state.State) {
	fmt.Println("Module informed about new state", startState.ToString())
}

func (client *AppStateClient) Start(state app_state.State) {
	fmt.Println("Client register as obserber and LOCK state", state.ToString())
	client.asm.RegisterObserver(client)
	client.asm.RegisterLockState(client, state)
}

func (client *AppStateClient) End(state app_state.State) {
	fmt.Println("Client unlock state", state.ToString())
	client.asm.UnlockState(client, state)
}

func main() {
	//dbHandler.RunHandler("test")

	asManager := app_state_manager.MakeAppStateManagerCtx()
	asClient := &AppStateClient{asManager, "A"}
	asClient.Start(app_state.LOADING)

	asClient2 := &AppStateClient{asManager, "B"}
	asClient2.Start(app_state.CONFIGURED)
	asManager.Start(app_state.INITIALIZING)

	asClient.End(app_state.LOADING)
	asClient2.End(app_state.CONFIGURED)

}
