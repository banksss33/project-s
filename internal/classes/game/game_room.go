package game

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/contrib/websocket"
)

type GameRoom struct {
	PlayerList  map[string]*Player // Key: userID from Discord, Value: *Player
	Setting     GameSetting
	State       GameState
	ActionEvent chan ActionEventType
	Broadcast   chan []byte
}

// create new game room Constructor
func NewGameRoom() *GameRoom {
	newRoom := &GameRoom{
		PlayerList:  make(map[string]*Player),
		Setting:     GameSetting{6, 420, 1},
		State:       GameState{},
		ActionEvent: make(chan ActionEventType),
		Broadcast:   make(chan []byte),
	}

	go newRoom.RunActionEventHandler()

	return newRoom
}

func (gr *GameRoom) RunActionEventHandler() {
	inUse := false
	for ActionEvent := range gr.ActionEvent {
		log.Println(ActionEvent)
		switch ActionEvent {
		case StartTimer:
			if !inUse {
				inUse = true
				go gr.StartTimer()
			}
		case CountDown:
			for _, player := range gr.PlayerList {
				timeLeft := fmt.Sprint(gr.State.TimeLeft)
				player.Conn.WriteMessage(websocket.TextMessage, []byte(timeLeft))
			}
		case ShutDown:
			close(gr.ActionEvent)
		}
	}
}

func (gr *GameRoom) StartTimer() {
	gr.State.TimeLeft = gr.Setting.Timer
	timer := time.NewTicker(time.Second)
	defer timer.Stop()

	for gr.State.TimeLeft > 0 {
		<-timer.C
		gr.State.TimeLeft--
		gr.ActionEvent <- CountDown
	}
}

// add new player to game room
func (gr *GameRoom) AddPlayer(userID string, conn *websocket.Conn) {
	gr.PlayerList[userID] = &Player{
		ConnectionStatus: StatusConnected,
		Conn:             conn,
		Score:            0,
	}
}

// use when player reconnected
func (gr *GameRoom) PlayerReconnected(userID string, conn *websocket.Conn) {
	reconnectedPlayer := gr.PlayerList[userID]
	reconnectedPlayer.ConnectionStatus = StatusConnected
	reconnectedPlayer.Conn = conn //assign new connection
}

// use when player disconnected
func (gr *GameRoom) PlayerDisconnected(userID string) {
	disconnectedPlayer := gr.PlayerList[userID]
	disconnectedPlayer.ConnectionStatus = StatusDisconnected
	if disconnectedPlayer.Conn != nil {
		disconnectedPlayer.Conn.Close()
		disconnectedPlayer.Conn = nil
	}
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
