package main

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type Player struct {
	ConnectionStatus string // CONNECTED | DISCONNECTED
	Conn             *websocket.Conn
	Score            uint // Score
}

// type GameSetting struct {
// 	RoundCount uint
// 	Timer      uint
// 	Spies      uint // 1 or 2
// }

type GameRoom struct {
	PlayerList map[string]*Player // Key: UserID from Discord Value: *Player
	Setting    string             // *GameSetting
	State      string             // *GameState
}

func createGameRoom() *GameRoom {
	return &GameRoom{
		PlayerList: make(map[string]*Player),
		Setting:    "",
		State:      "",
	}
}

func (gr *GameRoom) addPlayer(UserId string, InstanceId string, Conn *websocket.Conn) {
	gr.PlayerList[UserId] = &Player{
		ConnectionStatus: "CONNECTED",
		Conn:             Conn,
		Score:            0,
	}
}

func main() {
	App := fiber.New()

	App.Get("/ws", websocket.New(func(c *websocket.Conn) {
		var GameServer = make(map[string]*GameRoom) // Key: InstanceId from Discord Value: *GameRoom
		GameServer["QWE"] = createGameRoom()
	}))

	App.Listen(":8080")
}
