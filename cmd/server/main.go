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
		var newPlayer *game.Player = game.NewPlayer(userID, conn)

		mu.RLock()
		room, exists := gameServer[roomID]
		mu.RUnlock()

		if !exists {
			mu.Lock()
			room, exists := gameServer[roomID]
			if !exists {
				isClosedNotifier := make(chan bool)

				newRoom := game.NewGameRoom(isClosedNotifier)
				gameServer[roomID] = newRoom
				go func(id string) {
					<-isClosedNotifier

					//Room cleanup operation
					mu.Lock()
					closeRoom := gameServer[roomID]
					closeRoom.Cleanup()
					delete(gameServer, roomID)
					mu.Unlock()
				}(roomID)

				newRoom.PlayerRegister(newPlayer)
				mu.Unlock()
				return
			}

			room.PlayerRegister(newPlayer)
			mu.Unlock()
			return
		}

		room.PlayerRegister(newPlayer)
	}))

	App.Listen(":8080")
}
