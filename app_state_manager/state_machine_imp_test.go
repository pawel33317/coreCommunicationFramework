package app_state_manager

import (
	"fmt"
	"testing"

	"github.com/pawel33317/coreCommunicationFramework/app_state_manager/app_state"
)

//State machine client implementation
type AppStateTestClient struct {
	asm  AppStateClientHandler
	Name string
}

func (moduleData *AppStateTestClient) OnAppStateChanged(startState app_state.State) {
	fmt.Println("Module informed about new state", startState.ToString())
}
func (client *AppStateTestClient) RegisterClient() {
	client.asm.RegisterObserver(client)
}
func (client *AppStateTestClient) LockState(state app_state.State) {
	client.asm.RegisterLockState(client, state)
}
func (client *AppStateTestClient) UnlockState() {
	client.asm.UnlockState(client)
}

func TestMakeAppStateManagerImpInitialState(t *testing.T) {
	asManager := MakeAppStateManagerImp()
	if asManager.GetCurrentState() != app_state.DISABLED {
		t.Error("Incorrect initial state")
	}
}

func TestMakeAppStateManagerImpTargetState(t *testing.T) {
	asManager := MakeAppStateManagerImp()
	asManager.Start(app_state.INITIALIZING)
	if asManager.GetCurrentState() != app_state.OPERRABLE {
		t.Error("Incorrect target state")
	}
}

func TestMakeAppStateManagerImpRegisterClientAndBlockState(t *testing.T) {
	asManager := MakeAppStateManagerImp()
	smClient := &AppStateTestClient{asManager, "ClientA"}
	smClient.RegisterClient()
	smClient.LockState(app_state.LOADING)

	asManager.Start(app_state.INITIALIZING)

	//check if SM stop on blocked state
	if asManager.GetCurrentState() != app_state.LOADING {
		t.Error("Incorrect blocked state")
	}

	smClient.UnlockState()

	//check if SM go further after state unlock
	if asManager.GetCurrentState() != app_state.OPERRABLE {
		t.Error("Incorrect target state")
	}
}

func TestMakeAppStateManagerImpTwoClientsBlockStates(t *testing.T) {
	asManager := MakeAppStateManagerImp()
	smClientA := &AppStateTestClient{asManager, "ClientA"}
	smClientB := &AppStateTestClient{asManager, "ClientB"}

	smClientA.RegisterClient()
	smClientB.RegisterClient()

	smClientA.LockState(app_state.LOADING)
	smClientA.LockState(app_state.CONFIGURED)
	smClientB.LockState(app_state.LOADING)

	asManager.Start(app_state.INITIALIZING)

	//check if SM stop on first blocked state
	if asManager.GetCurrentState() != app_state.LOADING {
		t.Error("Incorrect blocked state")
	}

	smClientA.UnlockState()

	//check if SM still blocks state due to second client
	if asManager.GetCurrentState() != app_state.LOADING {
		t.Error("Incorrect blocked state")
	}

	smClientB.UnlockState()

	//check if SM unlocked first blocked state
	if asManager.GetCurrentState() != app_state.CONFIGURED {
		t.Error("Incorrect blocked state")
	}

	smClientA.UnlockState()

	//check if SM went to final state
	if asManager.GetCurrentState() != app_state.OPERRABLE {
		t.Error("Incorrect blocked state")
	}
}
