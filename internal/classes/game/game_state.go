package game

import "encoding/json"

type Action struct {
	ActionName    string
	ActionPayload json.RawMessage
}

type State interface {
	GetAction(Action) func(json.RawMessage)
	SetState(State)
}

// TimeLeft  int
// RoundLeft int
// Location  string
// Spies     int

type GameState struct {
	currentState State
}

func (gs *GameState) SetState(stateName State) {
	gs.currentState = stateName
}

type LobbyState struct {
	GameState
}

func (l *LobbyState) GetStateName() {

}

func (l *LobbyState) GetAction(p0 Action) func(json.RawMessage) {
	panic("TODO: Implement")
}
