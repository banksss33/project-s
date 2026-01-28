package types

import "encoding/json"

type Action struct {
	CallerID   string          `json:"caller_id"`
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

type AccusePayload struct {
	TargetUserID string `json:"target_user_id"`
}

type SpyGuessPayload struct {
	Location string `json:"location"`
}
