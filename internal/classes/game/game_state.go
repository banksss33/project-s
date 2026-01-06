package game

import (
	"encoding/json"
	"project-s/internal/types"
)

type State interface {
	GetStateName(types.ActionEvent)
}

// TimeLeft  int
// RoundLeft int
// Location  string
// Spies     int

type GameState struct {
	currentState State
}

func (gs *GameState) SetState(payload json.RawMessage) {
	State := map[string]State{
		"LOBBY": &LobbyState{},
	}

	var payloadData types.StatePayload
	json.Unmarshal(payload, &payloadData)

	gs.currentState = State[payloadData.SetState]
}

type LobbyState struct {
	GameState
}

func (l *LobbyState) GetStateName(event types.ActionEvent) {

}

func (l *LobbyState) GetPlayerList(event types.ActionEvent) {

}
