package main

import (
	"fmt"

	"github.com/pawel33317/coreCommunicationFramework/app_state_manager"
	"github.com/pawel33317/coreCommunicationFramework/app_state_manager/app_state"
)

type MyModule struct {
}

func (moduleData *MyModule) OnAppStateChanged(startState app_state.State) {
	fmt.Println("Module informed about new state", startState.ToString())
}

func main() {
	//dbHandler.RunHandler("test")

	asm := app_state_manager.MakeAppStateManagerData()
	aso := MyModule{}
	asm.RegisterObserver(&aso)
	asm.Start(app_state.INITIALIZING)
}
