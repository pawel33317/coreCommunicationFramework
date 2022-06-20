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
	INITIALIZED State = iota
	LOADED
	CONFIGURED
	OPERRABLE
	SHUTDOWN
	DISABLED
	ENABLED
)

//Converts SM state to string
func (state State) ToString() string {
	switch state {
	case INITIALIZED:
		return "INITIALIZED"
	case LOADED:
		return "LOADED"
	case CONFIGURED:
		return "CONFIGURED"
	case ENABLED:
		return "ENABLED"
	case SHUTDOWN:
		return "SHUTDOWN"
	case DISABLED:
		return "DISABLED"
	default:
		return "<UNKNOWN>"
	}
}

//Returns information whether current state is currently final
func (state State) IsTargetState() bool {
	switch state {
	case DISABLED, ENABLED:
		return true
	default:
		return false
	}
}

//Returns next state if current state is not target or error
func (state State) GetNextState() (State, error) {
	switch state {
	case INITIALIZED:
		return LOADED, nil
	case LOADED:
		return CONFIGURED, nil
	case CONFIGURED:
		return ENABLED, nil
	case ENABLED:
		return state, fmt.Errorf("Missing next state")
	case SHUTDOWN:
		return DISABLED, nil
	case DISABLED:
		return state, fmt.Errorf("Missing next state")
	default:
		return state, fmt.Errorf("Unknown state")
	}
}
