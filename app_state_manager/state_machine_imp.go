package app_state_manager

import (
	"fmt"

	"github.com/pawel33317/coreCommunicationFramework/app_state_manager/app_state"
	"github.com/pawel33317/coreCommunicationFramework/logger"
)

//private
type observerListUnique map[AppStateObserver]bool //set of AppStateObservers

//AppStateManagerImp structs
//implements AppStateManager
//implements AppStateClientHandler
type AppStateManagerImp struct {
	stateObserver                observerListUnique                     //set of AppStateObservers
	lockedState                  map[app_state.State]observerListUnique //set of states which contains set of AppStateObserver which are blocking particular state
	observerBlockingCurrentState observerListUnique                     //set of AppStateObservers which are blocking change of current state
	currentState                 app_state.State                        //current state
	logger                       logger.Logger
}

//MakeAppStateManagerImp constructor
func MakeAppStateManagerImp(smLogger logger.Logger) *AppStateManagerImp {
	return &AppStateManagerImp{
		stateObserver:                make(observerListUnique),
		lockedState:                  make(map[app_state.State]observerListUnique),
		observerBlockingCurrentState: make(observerListUnique),
		currentState:                 app_state.DISABLED,
		logger:                       smLogger,
	}
}

//private
func (asmData *AppStateManagerImp) informObservers() {
	asmData.logger.Log(logger.DEBUG, "SM", "Informing observers")
	for o := range asmData.stateObserver {
		o.OnAppStateChanged(asmData.currentState)
	}
}

//private
func (asmData *AppStateManagerImp) isCurrencStateBlocked() bool {
	return len(asmData.observerBlockingCurrentState) != 0
}

//private
func (asmData *AppStateManagerImp) processStates() {
	if asmData.isCurrencStateBlocked() {
		asmData.logger.Log(logger.INFO, "SM", "Current State is blocked by some observer")
		return
	}
	if asmData.currentState.IsTargetState() {
		asmData.logger.Log(logger.INFO, "SM", "Current State is target state, nothing to process", asmData.currentState.ToString())
		return
	}
	newState, err := asmData.currentState.GetNextState()
	if err != nil {
		asmData.logger.Log(logger.WARN, "SM", err.Error())
		return
	}
	asmData.changeState(newState)
	asmData.informObservers()
	asmData.processStates()
}

//private
func (asmData *AppStateManagerImp) changeState(newState app_state.State) {
	asmData.logger.Log(logger.DEBUG, "SM", "Changing state from:", asmData.currentState.ToString(), "to:", newState.ToString())
	asmData.currentState = newState
	asmData.setObserversBlockingCurrentState(newState)
}

//private
func (asmData *AppStateManagerImp) setObserversBlockingCurrentState(newState app_state.State) {
	if _, ok := asmData.lockedState[newState]; ok {
		asmData.observerBlockingCurrentState = make(observerListUnique)
		for k, v := range asmData.lockedState[newState] {
			asmData.observerBlockingCurrentState[k] = v
		}
	}
}

//from AppStateManager interface
func (asmData *AppStateManagerImp) Start(startState app_state.State) {
	asmData.changeState(startState)
	asmData.informObservers()
	asmData.processStates()
}

//from AppStateManager interface
func (asmData *AppStateManagerImp) GetCurrentState() app_state.State {
	return asmData.currentState
}

//from AppStateClientHandler interface
func (asmData *AppStateManagerImp) RegisterObserver(observer AppStateObserver) {
	asmData.stateObserver[observer] = true
}

//from AppStateClientHandler interface
func (asmData *AppStateManagerImp) RegisterLockState(observer AppStateObserver, state app_state.State) {
	asmData.logger.Log(logger.INFO, "SM", "Added lock state ", state.ToString(), " by application ", observer)
	if asmData.lockedState[state] == nil {
		asmData.lockedState[state] = make(observerListUnique)
	}
	asmData.lockedState[state][observer] = true
}

//from AppStateClientHandler interface
func (asmData *AppStateManagerImp) UnregisterLockState(observer AppStateObserver, state app_state.State) {
	asmData.logger.Log(logger.INFO, "SM", "Removed lock state ", state.ToString(), " by application ", observer)
	fmt.Println()
	delete(asmData.lockedState[state], observer)
}

//from AppStateClientHandler interface
func (asmData *AppStateManagerImp) UnlockState(observer AppStateObserver) {
	asmData.logger.Log(logger.INFO, "SM", "Unlocking state ", asmData.currentState.ToString(), " by application ", observer)
	delete(asmData.observerBlockingCurrentState, observer)
	asmData.processStates()
}
