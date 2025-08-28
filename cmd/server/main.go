package main

import (
	"log"
	"project-s/internal/classes/game"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func main() {
	App := fiber.New()
	var gameServer = make(map[string]*game.GameRoom)
	App.Get("/ws", websocket.New(func(conn *websocket.Conn) {
		roomID := conn.Query("roomID")
		userID := conn.Query("userID")

		if roomID == "" || userID == "" {
			closeMsg := websocket.FormatCloseMessage(
				websocket.CloseUnsupportedData,
				"Parameter required for roomID and userID",
			)
			conn.WriteMessage(websocket.CloseMessage, closeMsg)
			conn.Close()

			return
		}
		_, exists := gameServer[roomID]

		if !exists {
			gameServer[roomID] = game.NewGameRoom()
		}
		room := gameServer[roomID]
		room.AddPlayer(userID, conn)
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				break
			}
			log.Println(string(msg))

			room.ActionEvent <- 0
		}
	}))

	App.Listen(":8080")
}
