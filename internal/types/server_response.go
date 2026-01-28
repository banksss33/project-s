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
	Role     string `json:"role"`
	Location string `json:"location"`
}

type RolesRevealResponse struct {
	Winner      string        `json:"winner"`
	Location    string        `json:"location"`
	PlayersRole []PlayerRoles `json:"players_role"`
}

type PlayerRoles struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

type GameStatus struct {
	IsTimeRunning bool                         `json:"is_time_running"`
	IsVoting      bool                         `json:"is_voting"`
	IsRoundEnd    bool                         `json:"is_round_end"`
	RoundLeft     int                          `json:"round_left"`
	PlayerList    map[string]InGamePlayerStats `json:"player_list"`
}

type InGamePlayerStats struct {
	AlreadyVote bool `json:"already_vote"`
	Score       int  `json:"score"`
}
