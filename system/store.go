package system

import (
	"errors"
	"fmt"
)

// Action is an action
type Action struct {
	Type    interface{}
	Payload Payload
	Meta    interface{}
	Error   error
}

// Store is a store
type Store struct {
	State     State
	Logs      []Action
	Reducer   Reducer
	Listeners []Listener
}

// Reducer is a reducer function
type Reducer func(state State, action Action) State

// Payload is a payload
type Payload map[string]interface{}

// State is a state
type State map[string]interface{}

// Listener is a function that gets called everytime an action is dispatched
type Listener func()

// NewStore creates a new store
func NewStore(initial State, reducer Reducer) Store {
	s := Store{State: initial, Reducer: reducer}
	s.Dispatch(Action{Type: "@@STORE/INIT"})
	return s
}

// Dispatch dispatches action
func (s *Store) Dispatch(action Action) (Action, error) {
	if action.Type == nil {
		return action, errors.New("action must have a type")
	}

	for _, listen := range s.Listeners {
		listen()
	}

	s.Logs = append(s.Logs, action)
	s.State = s.Reducer(s.State, action)

	return action, nil
}

// Subscribe adds a listener
func (s *Store) Subscribe(listener Listener) func() {
	var isSubscribed = true

	s.Listeners = append(s.Listeners, listener)

	return func() {
		if !isSubscribed {
			return
		}

		index := 0

		for p, v := range s.Listeners {
			if &v == &listener {
				index = p
			}
		}

		s.Listeners = append(s.Listeners[:index], s.Listeners[index+1:]...)

		isSubscribed = false
	}
}

// LatestAction retrieves the latest action from the backlog
func (s *Store) LatestAction() Action {
	return s.Logs[len(s.Logs)-1]
}

// Debug logs the state tree and the action log.
func (s *Store) Debug() {
	fmt.Printf("[Store] State is %s, latest action is %s.\n", s.State, s.LatestAction().Type)
}

// MakeActionCreator creates an Action Creator
// func MakeActionCreator(input ...string) func(Payload | []string) Action {
// 	if len(input) > 1 {
// 		return func(payload Payload) Action {
// 			return Action{
// 				Type: input[1], Payload: Payload{},
// 			}
// 		}
// 	}
// 	return func(payload Payload) Action {
// 		return Action{Type: input[0], Payload: payload}
// 	}
// }
