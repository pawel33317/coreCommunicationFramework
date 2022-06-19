package app_state_manager

import (
	"testing"

	"github.com/pawel33317/coreCommunicationFramework/app_state_manager/app_state"
	"github.com/pawel33317/coreCommunicationFramework/logger"
)

//State machine client implementation
type AppStateTestClient struct {
	asm    AppStateClientHandler
	Name   string
	logger *logger.LogWrapper
}

func (client *AppStateTestClient) OnAppStateChanged(startState app_state.State) {
	client.logger.Log(logger.INFO, "Module informed about new state", startState.ToString())
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

func TestAppStateManagerImpInitialState(t *testing.T) {
	log := logger.NewLoggerImp()
	log.Disable()
	asManager := NewAppStateManagerImp(log)
	if asManager.GetCurrentState() != app_state.DISABLED {
		t.Error("Incorrect initial state")
	}
}

func TestAppStateManagerImpTargetState(t *testing.T) {
	log := logger.NewLoggerImp()
	log.Disable()
	asManager := NewAppStateManagerImp(log)
	asManager.Start(app_state.INITIALIZING)
	if asManager.GetCurrentState() != app_state.OPERRABLE {
		t.Error("Incorrect target state")
	}
}

func TestAppStateManagerImpRegisterClientAndBlockState(t *testing.T) {
	log := logger.NewLoggerImp()
	log.Disable()
	asManager := NewAppStateManagerImp(log)
	smClient := &AppStateTestClient{asManager, "ClientA", logger.NewLogWrapper(log, "TC")}
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

func TestAppStateManagerImpTwoClientsBlockStates(t *testing.T) {
	log := logger.NewLoggerImp()
	log.Disable()
	asManager := NewAppStateManagerImp(log)
	smClientA := &AppStateTestClient{asManager, "ClientA", logger.NewLogWrapper(log, "TC1")}
	smClientB := &AppStateTestClient{asManager, "ClientB", logger.NewLogWrapper(log, "TC2")}

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
