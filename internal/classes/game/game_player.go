package game

import "github.com/gofiber/contrib/websocket"

const StatusConnected = "CONNECTED"
const StatusDisconnected = "DISCONNECTED"

type Player struct {
	ConnectionStatus string // CONNECTED | DISCONNECTED
	Conn             *websocket.Conn
	Score            int // Score
}

func (p *Player) WritePump(message []byte) {
	p.Conn.WriteMessage(websocket.TextMessage, message)
}