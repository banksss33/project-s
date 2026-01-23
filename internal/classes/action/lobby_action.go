package action

import (
	"project-s/internal/classes/game"
	"project-s/internal/types"
)

func LobbyActionDispatcher() map[string]func(*game.GameRoom, *game.LobbyState, types.PlayerAction) {
	lobbyAction := make(map[string]func(*game.GameRoom, *game.LobbyState, types.PlayerAction))

	return lobbyAction
}

func joinGame(room *game.GameRoom, lobby *game.LobbyState, payload types.PlayerAction) {

}
