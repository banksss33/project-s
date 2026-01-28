package game

import (
	"project-s/internal/types"
	"sync"

	"github.com/gofiber/contrib/websocket"
)

type Player struct {
	UserID       string
	SendToPlayer chan types.ServerResponse
	IsConnected  bool // True = connect, False = notconnect
	Conn         *websocket.Conn
	mu           sync.Mutex
}

func NewPlayer(userID string, conn *websocket.Conn) *Player {
	return &Player{
		UserID:       userID,
		SendToPlayer: make(chan types.ServerResponse, 20),
		IsConnected:  true,
		Conn:         conn,
	}
}

func (p *Player) Reconnect(conn *websocket.Conn) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Conn = conn
	p.IsConnected = true
}

func (p *Player) disconnect() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.Conn != nil {
		p.IsConnected = false
		p.Conn.Close()
		p.Conn = nil
	}
}

func (p *Player) CreateReadPump(actionReceiver chan<- types.Action) {
	if p.Conn == nil {
		return
	}
	defer p.disconnect()

	var playerAction types.Action
	playerAction = types.Action{
		CallerID:   p.UserID,
		ActionName: "PLAYER_JOIN",
		Payload:    nil,
	}

	actionReceiver <- playerAction
	for {

		if err := p.Conn.ReadJSON(&playerAction); err != nil {
			playerAction = types.Action{
				CallerID:   p.UserID,
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
