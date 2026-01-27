package types

import (
	"encoding/json"
)

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

type PlayerRoleResponse struct {
	Role string `json:"role"`
}

type RolesRevealResponse struct {
	Winner      string        `json:"winner"`
	PlayersRole []PlayerRoles `json:"players_role"`
}

type PlayerRoles struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}
