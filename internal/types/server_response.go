package types

import "encoding/json"

type ServerResponse struct {
	ResponseName string          `json:"response_name"`
	Payload      json.RawMessage `json:"payload"`
}

type UpdateLobbyPlayerListResponse struct {
	Host       string   `json:"host"`
	Spectators []string `json:"spectators"`
	Players    []string `json:"players"`
}

type UpdateTimerResponse struct {
	TimerNow int `json:"timer_now"`
}
