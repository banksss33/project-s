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

	App.Get("/connect", websocket.New(func(conn *websocket.Conn) {
		roomID := conn.Query("roomID")
		userID := conn.Query("userID")

		mu.RLock()
		room, exists := gameServer[roomID]
		mu.RUnlock()

		if !exists {
			mu.Lock()
			room, exists = gameServer[roomID]
			if !exists {
				var newPlayer *game.Player = game.NewPlayer(userID, conn)
				isClosedNotifier := make(chan bool)

				newRoom := game.NewGameRoom(isClosedNotifier, newPlayer)
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

				go newPlayer.CreateWritePump()
				newPlayer.CreateReadPump(newRoom.ActionReceiver)
				return
			}
			mu.Unlock()
		}

		reconnectPlayer, playerExists := room.GetPlayerByID(userID)
		if playerExists {
			reconnectPlayer.Reconnect(conn)
			go reconnectPlayer.CreateWritePump()
			reconnectPlayer.CreateReadPump(room.ActionReceiver)
			return
		}

		var newPlayer *game.Player = game.NewPlayer(userID, conn)
		room.PlayerRegister(newPlayer)
		go newPlayer.CreateWritePump()
		newPlayer.CreateReadPump(room.ActionReceiver)
	}))

	App.Listen(":8080")
}
