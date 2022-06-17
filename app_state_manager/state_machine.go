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
	RegisterObserver(AppStateObserver)                     //registers to get inormation about system state changes
	RegisterLockState(AppStateObserver, app_state.State)   //state cannot be changed without confirmation "unlockState()" from AppStateObserver
	UnregisterLockState(AppStateObserver, app_state.State) //state can be changed without confirmation from AppStateObserver
	UnlockState(AppStateObserver, app_state.State)         //confirm that state can be changed
}

//AppStateManager implements this interface
type AppStateManager interface {
	Start(app_state.State)                    //start AppStateManager, tries to go from state INITIALIZING to OPERRABLE
	RequestStateChange(app_state.State) error //allows to move to acceptable state
	GetCurrentState() app_state.State         //returns current state
}

type observerListUnique map[AppStateObserver]bool //set of AppStateObservers

type AppStateManagerCtx struct {
	stateObserver                observerListUnique                     //set of AppStateObservers
	lockedState                  map[app_state.State]observerListUnique //set of states which contains set of AppStateObserver which are blocking particular state
	observerBlockingCurrentState observerListUnique                     //set of AppStateObservers which are blocking change of current state
	currentState                 app_state.State
	//TODO: logger
}

func MakeAppStateManagerCtx() *AppStateManagerCtx {
	return &AppStateManagerCtx{
		stateObserver:                make(observerListUnique),
		lockedState:                  make(map[app_state.State]observerListUnique),
		observerBlockingCurrentState: make(observerListUnique),
		currentState:                 0,
	}
}

func (asmData *AppStateManagerCtx) informObservers() {
	fmt.Println("Informing observers")
	for o := range asmData.stateObserver {
		o.OnAppStateChanged(asmData.currentState)
	}
}

func (asmData *AppStateManagerCtx) isCurrencStateBlocked() bool {
	return len(asmData.observerBlockingCurrentState) != 0
}

func (asmData *AppStateManagerCtx) processStates() {
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

func (asmData *AppStateManagerCtx) changeState(newState app_state.State) {
	fmt.Println("Changing state from:", asmData.currentState.ToString(), "to:", newState.ToString())
	asmData.currentState = newState
	if _, ok := asmData.lockedState[newState]; ok {
		asmData.observerBlockingCurrentState = make(observerListUnique)
		for k, v := range asmData.lockedState[newState] {
			asmData.observerBlockingCurrentState[k] = v
		}
	}
}

func (asmData *AppStateManagerCtx) Start(startState app_state.State) {
	asmData.currentState = startState

	fmt.Println("ASM started")

	if _, ok := asmData.lockedState[startState]; ok {
		asmData.observerBlockingCurrentState = make(observerListUnique)
		for k, v := range asmData.lockedState[startState] {
			asmData.observerBlockingCurrentState[k] = v
		}
	}
	asmData.informObservers()
	asmData.processStates()
}

func (asmData *AppStateManagerCtx) RegisterObserver(observer AppStateObserver) {
	asmData.stateObserver[observer] = true
}

func (asmData *AppStateManagerCtx) RegisterLockState(observer AppStateObserver, state app_state.State) {
	fmt.Println("Added lock state", state.ToString(), "by application", observer)
	if asmData.lockedState[state] == nil {
		asmData.lockedState[state] = make(observerListUnique)
	}
	asmData.lockedState[state][observer] = true
}

func (asmData *AppStateManagerCtx) UnregisterLockState(observer AppStateObserver, state app_state.State) {
	fmt.Println("Removed lock state", state.ToString(), "by application", observer)
	delete(asmData.lockedState[state], observer)
}

func (asmData *AppStateManagerCtx) UnlockState(observer AppStateObserver, state app_state.State) {
	fmt.Println("Unlocking state", state.ToString(), "by application", observer)
	delete(asmData.observerBlockingCurrentState, observer)
	asmData.processStates()
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
