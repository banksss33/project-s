package types

import (
	"encoding/json"
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
	Score int
	Roles string
}

type GameSetting struct {
	Spies     int
	Timer     int
	Locations map[string][]string //key: location name| value: location roles
}
