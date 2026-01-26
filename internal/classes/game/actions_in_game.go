package game

import (
	"project-s/internal/types"
)

func InGameActionDispatcher() map[string]func(*GameRoom, *InGameState, types.PlayerAction) {
	InGameAction := make(map[string]func(*GameRoom, *InGameState, types.PlayerAction))

	return InGameAction
}

func NewGame(room *GameRoom, inGameState *InGameState, playerAction types.PlayerAction) {
	// var gameStartPayload types.GameStartPayload
	// json.Unmarshal(playerAction.Payload, &gameStartPayload)
	// gameSetting := types.GameSetting{
	// 	Round:     gameStartPayload.GameSetting.Round,
	// 	Spies:     gameStartPayload.GameSetting.Spies,
	// 	Timer:     gameStartPayload.GameSetting.Timer,
	// 	Locations: gameStartPayload.GameSetting.Locations,
	// }

	// var playerList []*Player
	// var spectatorList []*Player

	// for _, userID := range gameStartPayload.Players {
	// 	player := room.GetPlayerByID(userID)

	// }
	// inGameState.Init(gameSetting)
}
