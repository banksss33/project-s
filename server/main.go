package main

import (
	"project-s/game"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func main() {
	App := fiber.New()
	var gameServer = make(map[string]*game.GameRoom)
	App.Get("/ws", websocket.New(func(conn *websocket.Conn) {
		roomID := conn.Query("roomID")
		userID := conn.Query("userID")

		_, exists := gameServer[roomID]
		if !exists {
			gameServer[roomID] = game.NewGameRoom()
		}
		room := gameServer[roomID]
		room.AddPlayer(userID, conn)

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	}))

	App.Listen(":8080")
}
