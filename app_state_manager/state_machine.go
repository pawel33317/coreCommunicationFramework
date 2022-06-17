package app_state_manager

import (
	"fmt"

	"github.com/pawel33317/coreCommunicationFramework/app_state_manager/app_state"
)

//Receives information abour application state change
type AppStateObserver interface {
	OnAppStateChanged(app_state.State)
}

//AppStateManager implements this interface
type AppStateClient interface {
	RegisterObserver(AppStateObserver)                      //registers to get inormation about system state changes
	RegisterLockState(*AppStateObserver, app_state.State)   //state cannot be changed without confirmation "unlockState()" from AppStateObserver
	UnregisterLockState(*AppStateObserver, app_state.State) //state can be changed without confirmation from AppStateObserver
	UnlockState(*AppStateObserver, app_state.State)         //confirm that state can be changed
}

//AppStateManager implements this interface
type AppStateManager interface {
	Start(app_state.State)                    //start AppStateManager, tries to go from state INITIALIZING to OPERRABLE
	RequestStateChange(app_state.State) error //allows to move to acceptable state
	GetCurrentState() app_state.State         //returns current state
}

type observerListUnique map[AppStateObserver]bool //set of AppStateObservers

type AppStateManagerData struct {
	stateObserver                observerListUnique                     //set of AppStateObservers
	lockedState                  map[app_state.State]observerListUnique //set of states which contains set of AppStateObserver which are blocking particular state
	observerBlockingCurrentState observerListUnique                     //set of AppStateObservers which are blocking change of current state
	currentState                 app_state.State
	//TODO: logger
}

func MakeAppStateManagerData() *AppStateManagerData {
	return &AppStateManagerData{
		stateObserver:                make(observerListUnique),
		lockedState:                  make(map[app_state.State]observerListUnique),
		observerBlockingCurrentState: make(observerListUnique),
		currentState:                 0,
	}
}

func (asmData *AppStateManagerData) informObservers() {
	fmt.Println("Informing observers")
	for o := range asmData.stateObserver {
		o.OnAppStateChanged(asmData.currentState)
	}
}

func (asmData *AppStateManagerData) isCurrencStateBlocked() bool {
	return len(asmData.observerBlockingCurrentState) != 0
}

func (asmData *AppStateManagerData) processStates() {
	if asmData.isCurrencStateBlocked() {
		fmt.Println("Current State is blocked be some observer")
		return
	}
	if asmData.currentState.IsTargetState() {
		fmt.Println("Current State target state, nothing to process")
		return
	}
	newState, err := asmData.currentState.GetNextState()
	if err != nil {
		fmt.Println(err)
		return
	}
	asmData.changeState(newState)
	asmData.informObservers()
	asmData.processStates()
}

func (asmData *AppStateManagerData) changeState(newState app_state.State) {
	fmt.Println("Changing state from:", asmData.currentState.ToString(), "to:", newState.ToString())
	asmData.currentState = newState
	if _, ok := asmData.lockedState[newState]; ok {
		asmData.observerBlockingCurrentState = asmData.lockedState[newState]
	}
}

func (asmData *AppStateManagerData) Start(startState app_state.State) {
	asmData.currentState = startState

	fmt.Println("ASM started")

	if _, ok := asmData.lockedState[startState]; ok {
		asmData.observerBlockingCurrentState = asmData.lockedState[startState]
	}
	asmData.informObservers()
	asmData.processStates()
}

func (asmData *AppStateManagerData) RegisterObserver(observer AppStateObserver) {
	asmData.stateObserver[observer] = true
}

/*
stateObserver
	OBSERVER_1
	OBSERVER_2

lockedState
	STATE_1
		OBSERVER_1
		OBSERVER_2
	STATE_2
		OBSERVER_1
		OBSERVER_3

observersBlockingCurrentState
	OBSERVER_1
	OBSERVER_2
*/
