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

func (gr *GameRoom) EventHandler() {
	for event := range gr.EventReceiver {
		go gr.Publish.Notify(event.ActionName, event)
	}
}

// add new player to game room
func (gr *GameRoom) AddPlayer(userID string, conn *websocket.Conn) {
	var player *Player = NewPlayer(userID, conn)
	go player.Read(gr.EventReceiver)
	gr.PlayerList[player.UserID] = player
}

// use when player reconnected
func (gr *GameRoom) PlayerReconnected(userID string, newConn *websocket.Conn) {
	var reconnectedPlayer *Player = gr.PlayerList[userID]
	go reconnectedPlayer.Read(gr.EventReceiver)
	reconnectedPlayer.Reconnect(newConn)
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
