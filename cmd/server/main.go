package main

import (
	"project-s/internal/classes/game"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func main() {
	App := fiber.New()
	var mu sync.RWMutex
	var gameServer = make(map[string]*game.GameRoom)

	App.Get("/ws", websocket.New(func(conn *websocket.Conn) {
		roomID := conn.Query("roomID")
		userID := conn.Query("userID")
		var player *game.Player = game.NewPlayer(userID, conn)

		mu.Lock()
		if _, exists := gameServer[roomID]; !exists {
			isClose := make(chan string)
			gameServer[roomID] = game.NewGameRoom(isClose)

			go func() {
				id := <-isClose
				delete(gameServer, id)
			}()
		}
		mu.Unlock()

		mu.RLock()
		gameServer[roomID].PlayerRegister(player)
		mu.RUnlock()
	}))

	App.Listen(":8080")
}
