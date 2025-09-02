package game

import (
	"encoding/json"
	"project-s/internal/types"

	"github.com/gofiber/contrib/websocket"
)

const StatusConnected = "CONNECTED"
const StatusDisconnected = "DISCONNECTED"

type Player struct {
	UserID           string
	Roles            string
	ConnectionStatus string // CONNECTED | DISCONNECTED
	Conn             *websocket.Conn
	Score            int // Score
}

func NewPlayer(userID string, conn *websocket.Conn) *Player {
	return &Player{
		UserID:           userID,
		ConnectionStatus: StatusConnected,
		Conn:             conn,
		Score:            0,
	}
}

func (p *Player) Reconnect(reconnect *websocket.Conn) {
	p.ConnectionStatus = "CONNECTED"
	p.Conn = reconnect
}

func (p *Player) Disconnect() {
	p.ConnectionStatus = "DISCONNECTED"
	if p.Conn != nil {
		p.Conn.Close()
		p.Conn = nil
	}
}

func (p *Player) Listener(ReceiveEvent chan types.ActionEvent) {
	defer p.Disconnect()
	for {
		_, msg, err := p.Conn.ReadMessage()
		if err != nil {
			break
		}

		var clientAction types.ClientAction
		json.Unmarshal(msg, &clientAction)
		ActionEvent := types.ActionEvent{
			PlayerConn:   p.Conn,
			ClientAction: clientAction,
		}

		ReceiveEvent <- ActionEvent

	}
}
