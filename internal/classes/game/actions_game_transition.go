package game

import (
	"project-s/internal/types"
)

func GameTransitionActionDispatcher() map[string]func(*GameRoom, *InGameState, *LobbyState, types.Action) {
	transitionAction := make(map[string]func(*GameRoom, *InGameState, *LobbyState, types.Action))

	transitionAction["NEW_ROUND"] = startNewRound
	transitionAction["BACK_TO_LOBBY"] = backToLobby

	return transitionAction
}

func backToLobby(room *GameRoom, inGame *InGameState, lobby *LobbyState, payload types.Action) {
	inGame = NewInGameState()
	room.state = "LOBBY_STATE"
}

func startNewRound(room *GameRoom, inGame *InGameState, lobby *LobbyState, payload types.Action) {
	room.state = "IN_GAME_STATE"
	inGame.NewRound()
	inGame.StartNewTimer(room.countdown, room.timeUp)
	SendPlayerRolesAndLocation(inGame)
	broadcastGameState(room, inGame)
}
