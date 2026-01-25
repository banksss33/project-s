package types

type UpdateLobbyPlayerListResponse struct {
	Host       string   `json:"host"`
	Spectators []string `json:"spectators"`
	Players    []string `json:"players"`
}
