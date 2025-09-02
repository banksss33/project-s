package game

import (
	"project-s/internal/services/pubsub"
	"project-s/internal/types"

	"github.com/gofiber/contrib/websocket"
)

type GameRoom struct {
	PlayerList   map[string]*Player // Key: userID from Discord, Value: *Player
	Publish      *pubsub.Publisher
	ReceiveEvent chan types.ActionEvent
	GameState
}

// create new game room Constructor
func NewGameRoom() *GameRoom {
	newRoom := &GameRoom{
		PlayerList:   make(map[string]*Player),
		Publish:      pubsub.NewPublisher(),
		ReceiveEvent: make(chan types.ActionEvent, 5),
	}

	return newRoom
}

func (gr *GameRoom) EventHandler() {
	for e := range gr.ReceiveEvent {
		go gr.Publish.Notify(e.ActionName, e)
	}
}

// add new player to game room
func (gr *GameRoom) AddPlayer(newPlayer *Player) {
	gr.PlayerList[newPlayer.UserID] = newPlayer
}

// use when player reconnected
func (gr *GameRoom) PlayerReconnected(userID string, newConn *websocket.Conn) {
	reconnectedPlayer := gr.PlayerList[userID]
	reconnectedPlayer.Reconnect(newConn)
}

// use when player disconnected
func (gr *GameRoom) PlayerDisconnected(userID string) {
	disconnectedPlayer := gr.PlayerList[userID]
	disconnectedPlayer.Disconnect()
}

// If all player disconnected return true
func (gr *GameRoom) IsRoomEmpty() bool {
	for _, player := range gr.PlayerList {
		if player.ConnectionStatus == StatusConnected {
			return false
		}
	}
	return true
}
