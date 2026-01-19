package game

import (
	"project-s/internal/types"

	"github.com/gofiber/contrib/websocket"
)

type Player struct {
	UserID      string
	SendJSON    chan types.ServerResponse
	IsConnected bool // True = connect, False = notconnect
	Conn        *websocket.Conn
}

func NewPlayer(userID string, conn *websocket.Conn) *Player {
	return &Player{
		UserID:      userID,
		SendJSON:    make(chan types.ServerResponse),
		IsConnected: true,
		Conn:        conn,
	}
}

func (p *Player) disconnect() {
	p.IsConnected = false
	if p.Conn != nil {
		p.Conn.Close()
		p.Conn = nil
	}
}

func (p *Player) CreateReadPump(actionReceiver chan<- types.PlayerAction) {
	defer p.disconnect()
	for {
		var playerAction types.PlayerAction

		if err := p.Conn.ReadJSON(&playerAction); err != nil {
			playerAction = types.PlayerAction{
				UserID:     p.UserID,
				ActionName: "PLAYER_DISCONNECT",
				Payload:    nil,
			}
			actionReceiver <- playerAction
			break
		}

		actionReceiver <- playerAction
	}
}

func (p *Player) CreateWritePump() {
	defer p.disconnect()
	for JSON := range p.SendJSON {
		err := p.Conn.WriteJSON(JSON)
		if err != nil {
			break
		}
	}
}
