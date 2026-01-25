package types

import (
	"encoding/json"
	"time"
)

type PlayerAction struct {
	UserID     string          `json:"user_id"`
	ActionName string          `json:"action_name"`
	Payload    json.RawMessage `json:"payload"`
}

type ServerResponse struct {
	ResponseName string          `json:"response_name"`
	Payload      json.RawMessage `json:"payload"`
}

type PlayerStatus struct {
	Score       int
	Roles       string
	AlreadyVote bool
}

type GameSetting struct {
	Round     int                 `json:"round"`
	Spies     int                 `json:"spies"`
	Timer     int                 `json:"timer"`
	Locations map[string][]string `json:"locations"` //key: location name| value: location roles
}

type GameTimer struct {
	Tick      *time.Ticker
	IsRunning bool
	Countdown int
}

type GameStatus struct {
	IsTimeRunning bool
	IsVoting      bool
	IsRoundEnd    bool
	RoundLeft     int
	PlayerStats   map[string]PlayerStatus
	Spectator     []string
}
