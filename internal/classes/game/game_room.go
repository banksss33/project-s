package game

import (
	"project-s/internal/types"
)

type GameRoom struct {
	playerList     map[string]*Player // Key: userID from Discord, Value: *Player
	state          string
	broadcast      chan types.ServerResponse
	actionReceiver chan types.PlayerAction

	gameClose chan bool
}

// create new game room Constructor
func NewGameRoom(gameClosedNotifier chan bool, host *Player) *GameRoom {
	newRoom := &GameRoom{
		playerList:     make(map[string]*Player),
		state:          "LOBBY_STATE",
		broadcast:      make(chan types.ServerResponse, 5),
		actionReceiver: make(chan types.PlayerAction, 5),
		gameClose:      gameClosedNotifier,
	}

	go newRoom.broadcastInit()
	go newRoom.actionProcessorInit()

	gameCreatedAction := types.PlayerAction{
		UserID:     host.UserID,
		ActionName: "GAME_CREATED",
		Payload:    nil,
	}

	newRoom.actionReceiver <- gameCreatedAction

	return newRoom
}

func (gr *GameRoom) broadcastInit() {
	for broadcastItem := range gr.broadcast {
		for _, player := range gr.playerList {
			player.SendToPlayer <- broadcastItem
		}
	}
}

func (gr *GameRoom) actionProcessorInit() {
	// var lobby *LobbyState
	// var inGame *InGameState

	// for action := range gr.actionReceiver {
	// 	switch gr.state {
	// 	case "LOBBY_STATE":

	// 	}
	// }
}

// add new player to game room
func (gr *GameRoom) PlayerRegister(player *Player) {
	gr.playerList[player.UserID] = player

	go player.CreateReadPump(gr.actionReceiver)
	go player.CreateWritePump()
}

func (gr *GameRoom) Cleanup() {
	if !gr.isEmpty() {
		return
	}

	close(gr.actionReceiver)
	close(gr.broadcast)
}

//#region action method - helper

// If all player disconnected return true
func (gr *GameRoom) isEmpty() bool {
	for _, player := range gr.playerList {
		if player.IsConnected {
			return false
		}
	}
	return true
}
