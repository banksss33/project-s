package game

import (
	"encoding/json"
	"project-s/internal/types"

	"github.com/gofiber/contrib/websocket"
)

type Player struct {
	UserID      string
	SendJSON    chan json.RawMessage
	IsConnected bool // True = connect, False = notconnect
	Conn        *websocket.Conn
}

func NewPlayer(userID string, conn *websocket.Conn) *Player {
	return &Player{
		UserID:      userID,
		SendJSON:    make(chan json.RawMessage),
		IsConnected: true,
		Conn:        conn,
	}
}

func (p *Player) Reconnect(conn *websocket.Conn) {
	p.IsConnected = true
	p.Conn = conn
}

func (p *Player) Disconnect() {
	p.IsConnected = false
	if p.Conn != nil {
		p.Conn.Close()
		p.Conn = nil
	}
}

func (p *Player) CreateReadPump(eventReceiver chan<- types.PlayerAction) {
	defer p.Disconnect()
	for {
		var playerAction types.PlayerAction

		if err := p.Conn.ReadJSON(&playerAction); err != nil {
			break
		}

		eventReceiver <- playerAction
	}
}

func (p *Player) CreateWritePump() {
	defer p.Disconnect()
	for JSON := range p.SendJSON {
		err := p.Conn.WriteJSON(JSON)
		if err != nil {
			break
		}
	}
}
