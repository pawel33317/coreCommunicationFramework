package app_state_manager

import (
	"github.com/pawel33317/coreCommunicationFramework/app_state_manager/app_state"
)

/*
 *  AppStateObserver interface interface
 *  AppStateClientHandler interface
 *  AppStateManager interface
 */

//Receives information about application state change
type AppStateObserver interface {
	OnAppStateChanged(app_state.State)
}

//AppStateManagerImp implements this interface
type AppStateClientHandler interface {
	RegisterObserver(AppStateObserver)                     //registers to get inormation about system state changes
	RegisterLockState(AppStateObserver, app_state.State)   //state cannot be changed without confirmation "unlockState()" from AppStateObserver
	UnregisterLockState(AppStateObserver, app_state.State) //state can be changed without confirmation from AppStateObserver
	UnlockState(AppStateObserver)                          //confirm that state can be changed
	GetCurrentState() app_state.State                      //returns current state
}

//AppStateManagerImp implements this interface
type AppStateManager interface {
	Start(app_state.State)                    //start AppStateManager, tries to go from state INITIALIZING to OPERRABLE
	RequestStateChange(app_state.State) error //allows to move to acceptable state, bool force
}
