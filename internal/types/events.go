package types

import (
	"encoding/json"

	"github.com/gofiber/contrib/websocket"
)

type ClientAction struct {
	ActionName    string          `json:"actionName"`
	ActionPayload json.RawMessage `json:"actionPayload"`
}

type ActionEvent struct {
	PlayerConn *websocket.Conn
	ClientAction
}
