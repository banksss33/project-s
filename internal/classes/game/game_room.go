package game

import (
	"project-s/internal/types"
	"sync"
)

type GameRoom struct {
	state          string
	PlayerList     map[string]*Player // Key: userID from Discord, Value: *Player
	Broadcast      chan types.ServerResponse
	ActionReceiver chan types.PlayerAction
	gameClose      chan bool
	mu             sync.Mutex
}

// create new game room Constructor
func NewGameRoom(gameClose chan bool, host *Player) *GameRoom {
	newRoom := &GameRoom{
		PlayerList:     make(map[string]*Player),
		state:          "LOBBY_STATE",
		Broadcast:      make(chan types.ServerResponse, 5),
		ActionReceiver: make(chan types.PlayerAction, 5),
		gameClose:      gameClose,
	}

	go newRoom.broadcastInit()
	go newRoom.actionProcessorInit()

	return newRoom
}

func (gr *GameRoom) broadcastInit() {
	for broadcastItem := range gr.Broadcast {
		for _, player := range gr.PlayerList {
			player.SendToPlayer <- broadcastItem
		}
	}
}

func (gr *GameRoom) actionProcessorInit() {
	var lobby *LobbyState = &LobbyState{}
	// var inGame *InGameState

	lobbyAction := LobbyActionDispatcher()
	for playerAction := range gr.ActionReceiver {
		switch gr.state {
		case "LOBBY_STATE":
			action := lobbyAction[playerAction.ActionName]
			action(gr, lobby, playerAction)
		}
	}
}

// add new player to game room
func (gr *GameRoom) PlayerRegister(player *Player) {
	gr.mu.Lock()
	defer gr.mu.Unlock()

	gr.PlayerList[player.UserID] = player
}

func (gr *GameRoom) Cleanup() {
	if !gr.isEmpty() {
		return
	}

	close(gr.ActionReceiver)
	close(gr.Broadcast)
}

//#region action method - helper

// If all player disconnected return true
func (gr *GameRoom) isEmpty() bool {
	for _, player := range gr.PlayerList {
		if player.IsConnected {
			return false
		}
	}
	return true
}
