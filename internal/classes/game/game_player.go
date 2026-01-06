package game

import (
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
	p.ConnectionStatus = StatusConnected
	p.Conn = reconnect
}

func (p *Player) Disconnect() {
	p.ConnectionStatus = StatusDisconnected
	if p.Conn != nil {
		p.Conn.Close()
		p.Conn = nil
	}
}

func (p *Player) Read(eventReceiver chan<- types.ActionEvent) {
	defer p.Disconnect()
	for {
		var clientAction types.ActionEvent

		if err := p.Conn.ReadJSON(&clientAction); err != nil {
			break
		}

		eventReceiver <- clientAction
	}
}

func (p *Player) Send() {

}
