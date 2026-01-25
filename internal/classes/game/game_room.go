package game

import (
	"project-s/internal/types"
	"sync"
)

type GameRoom struct {
	state          string
	playerList     map[string]*Player // Key: userID from Discord, Value: *Player
	broadcast      chan types.ServerResponse
	ActionReceiver chan types.PlayerAction
	gameClose      chan bool
	mu             sync.Mutex
}

// create new game room Constructor
func NewGameRoom(gameClose chan bool, host *Player) *GameRoom {
	newRoom := &GameRoom{
		playerList:     make(map[string]*Player),
		state:          "LOBBY_STATE",
		broadcast:      make(chan types.ServerResponse, 5),
		ActionReceiver: make(chan types.PlayerAction, 5),
		gameClose:      gameClose,
	}

	go newRoom.broadcastInit()
	go newRoom.actionProcessorInit()

	gameCreatedAction := types.PlayerAction{
		UserID:     host.UserID,
		ActionName: "GAME_CREATED",
		Payload:    nil,
	}

	newRoom.ActionReceiver <- gameCreatedAction

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
	gr.mu.Lock()
	defer gr.mu.Unlock()

	gr.playerList[player.UserID] = player
}

func (gr *GameRoom) Cleanup() {
	if !gr.isEmpty() {
		return
	}

	close(gr.ActionReceiver)
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
