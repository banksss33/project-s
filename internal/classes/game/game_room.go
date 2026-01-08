package game

import (
	"project-s/internal/services/pubsub"
	"project-s/internal/types"

	"github.com/gofiber/contrib/websocket"
)

type GameRoom struct {
	PlayerList    map[string]*Player // Key: userID from Discord, Value: *Player
	Publish       *pubsub.Publisher
	EventReceiver chan types.ActionEvent
	GameState
}

// create new game room Constructor
func NewGameRoom() *GameRoom {
	newRoom := &GameRoom{
		PlayerList:    make(map[string]*Player),
		Publish:       pubsub.NewPublisher(),
		EventReceiver: make(chan types.ActionEvent, 5),
	}

	return newRoom
}

// add new player to game room
func (gr *GameRoom) AddPlayer(userID string, conn *websocket.Conn) {
	var player *Player = NewPlayer(userID, conn)
	gr.PlayerList[player.UserID] = player

	go player.CreateReadPump(gr.EventReceiver)
	go player.CreateWritePump()
}

// use when player reconnected
func (gr *GameRoom) PlayerReconnected(userID string, newConn *websocket.Conn) {
	var reconnectedPlayer *Player = gr.PlayerList[userID]
	reconnectedPlayer.Reconnect(newConn)

	go reconnectedPlayer.CreateReadPump(gr.EventReceiver)
	go reconnectedPlayer.CreateWritePump()
}

// use when player disconnected
func (gr *GameRoom) PlayerDisconnected(userID string) {
	var disconnectedPlayer *Player = gr.PlayerList[userID]
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
