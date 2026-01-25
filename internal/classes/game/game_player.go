package game

import (
	"project-s/internal/types"

	"github.com/gofiber/contrib/websocket"
)

type Player struct {
	UserID       string
	SendToPlayer chan types.ServerResponse
	IsConnected  bool // True = connect, False = notconnect
	Conn         *websocket.Conn
}

func NewPlayer(userID string, conn *websocket.Conn) *Player {
	return &Player{
		UserID:       userID,
		SendToPlayer: make(chan types.ServerResponse),
		IsConnected:  true,
		Conn:         conn,
	}
}

func (p *Player) disconnect() {
	if p.Conn != nil {
		p.IsConnected = false
		p.Conn.Close()
		p.Conn = nil
		close(p.SendToPlayer)
	}
}

func (p *Player) CreateReadPump(actionReceiver chan<- types.PlayerAction) {
	if p.Conn == nil {
		return
	}
	defer p.disconnect()

	var playerAction types.PlayerAction
	playerAction = types.PlayerAction{
		UserID:     p.UserID,
		ActionName: "PLAYER_JOIN",
		Payload:    nil,
	}

	actionReceiver <- playerAction
	for {

		if err := p.Conn.ReadJSON(&playerAction); err != nil {
			playerAction = types.PlayerAction{
				UserID:     p.UserID,
				ActionName: "PLAYER_LEFT",
				Payload:    nil,
			}
			actionReceiver <- playerAction
			break
		}

		actionReceiver <- playerAction
	}
}

func (p *Player) CreateWritePump() {
	if p.Conn == nil {
		return
	}
	for JSON := range p.SendToPlayer {
		err := p.Conn.WriteJSON(JSON)
		if err != nil {
			break
		}
	}
}
