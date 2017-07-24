package system

// ActionType defines the type used in action's type field.
type ActionType string

const (
	// ONLINE is an action type
	ONLINE ActionType = "ONLINE"
	// OFFLINE is an action type
	OFFLINE ActionType = "OFFLINE"
)

// OfflineReducer is a reducer
func OfflineReducer(state State, action Action) State {
	switch action.Type {
	case OFFLINE:
		state["status"] = "App Offline"
	case ONLINE:
		state["status"] = "App Online"
	}
	return state
}

func setOnline() Action {
	return Action{Type: ONLINE}
}

func setOffline() Action {
	return Action{Type: OFFLINE}
}

// OfflineHandler handles offline actions
func OfflineHandler() Store {
	ost := NewStore(State{"status": "App Offline"}, OfflineReducer)

	unsub := ost.Subscribe(func() {
		ost.Debug()
	})

	ost.Dispatch(setOnline())
	ost.Dispatch(setOffline())

	ost.Dispatch(Action{Type: ONLINE})
	unsub()

	return ost
}
