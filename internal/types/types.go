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
	ResponseName string          `json:"json:response_name"`
	Payload      json.RawMessage `json:"payload"`
}

type PlayerStatus struct {
	Score       int
	Roles       string
	AlreadyVote bool
}

type GameSetting struct {
	Round     int
	Spies     int
	Timer     int
	Locations map[string][]string //key: location name| value: location roles
}

type GameTimer struct {
	Tick      *time.Ticker
	IsRunning bool
	Countdown int
}
