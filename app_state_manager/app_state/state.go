package app_state

import "fmt"

/*
 * State states
 *   - 11 SM states
 *   - ToString()
 *   - IsTargetState()
 *   - GetNextState()
 */

//SM state
type State int

//List of SM states
const (
	INITIALIZING State = iota
	INITIALIZED
	LOADING
	LOADED
	CONFIGURING
	CONFIGURED
	OPERRABLE
	SHUTTINGDOWN
	TURNEDOFF
	DISABLED
	ENABLED
)

//Converts SM state to string
func (state State) ToString() string {
	switch state {
	case INITIALIZING:
		return "INITIALIZING"
	case INITIALIZED:
		return "INITIALIZED"
	case LOADING:
		return "LOADING"
	case LOADED:
		return "LOADED"
	case CONFIGURING:
		return "CONFIGURING"
	case CONFIGURED:
		return "CONFIGURED"
	case OPERRABLE:
		return "OPERRABLE"
	case SHUTTINGDOWN:
		return "SHUTTINGDOWN"
	case TURNEDOFF:
		return "TURNEDOFF"
	case DISABLED:
		return "DISABLED"
	case ENABLED:
		return "ENABLED"
	default:
		return "<UNKNOWN>"
	}
}

//Returns information whether current state is currently final
func (state State) IsTargetState() bool {
	switch state {
	case OPERRABLE, TURNEDOFF, DISABLED, ENABLED:
		return true
	default:
		return false
	}
}

//Returns next state if current state is not target or error
func (state State) GetNextState() (State, error) {
	switch state {
	case INITIALIZING:
		return INITIALIZED, nil
	case INITIALIZED:
		return LOADING, nil
	case LOADING:
		return LOADED, nil
	case LOADED:
		return CONFIGURING, nil
	case CONFIGURING:
		return CONFIGURED, nil
	case CONFIGURED:
		return OPERRABLE, nil
	case OPERRABLE:
		return state, fmt.Errorf("Missing next state")
	case SHUTTINGDOWN:
		return TURNEDOFF, nil
	case TURNEDOFF:
		return state, fmt.Errorf("Missing next state")
	case DISABLED:
		return state, fmt.Errorf("Missing next state")
	case ENABLED:
		return state, fmt.Errorf("Missing next state")
	default:
		return state, fmt.Errorf("Unknown state")
	}
}
