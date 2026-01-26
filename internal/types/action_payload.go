package types

import "encoding/json"

type PlayerAction struct {
	UserID     string          `json:"user_id"`
	ActionName string          `json:"action_name"`
	Payload    json.RawMessage `json:"payload"`
}

type GameStartPayload struct {
	Players     []string           `json:"players"`
	Spectators  []string           `json:"spectators"`
	GameSetting GameSettingPayload `json:"game_setting"`
}

type GameSettingPayload struct {
	Round     int                 `json:"round"`
	Spies     int                 `json:"spies"`
	Timer     int                 `json:"timer"`
	Locations map[string][]string `json:"locations"` //key: location name| value: location roles
}
