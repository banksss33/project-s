// package game_test

// import (
// 	"project-s/internal/classes/game"
// 	"testing"
// )

// func TestIsRoomEmpty(t *testing.T) {
// 	t.Run("newly created room should be empty", func(t *testing.T) {
// 		room := game.NewGameRoom("TEST")

// 		if !room.IsEmpty() {
// 			t.Error("Expected new room to be empty")
// 		}
// 	})

// 	t.Run("room with one connected player should not be empty", func(t *testing.T) {
// 		room := game.NewGameRoom("TEST")
// 		newPlayer := game.NewPlayer("p1", nil)
// 		room.PlayerRegister(newPlayer)

// 		if room.IsEmpty() {
// 			t.Error("Expected room with one player to not be empty")
// 		}
// 	})

// 	t.Run("room after all players disconnect should be empty", func(t *testing.T) {
// 		room := game.NewGameRoom("TEST")
// 		newPlayer := game.NewPlayer("p1", nil)
// 		room.PlayerRegister(newPlayer)
// 		room.RemovePlayer(newPlayer)

// 		if !room.IsEmpty() {
// 			t.Error("Expected room to be empty after player disconnected")
// 		}
// 	})

// 	t.Run("room with two player should not be empty", func(t *testing.T) {
// 		room := game.NewGameRoom("TEST")

// 		p1 := game.NewPlayer("p1", nil)
// 		p2 := game.NewPlayer("p2", nil)
// 		room.PlayerRegister(p1)
// 		room.PlayerRegister(p2)

// 		if room.IsEmpty() {
// 			t.Error("Expected room with two players to not be empty")
// 		}
// 	})

// 	t.Run("room with one player after player disconnect should not be empty", func(t *testing.T) {
// 		room := game.NewGameRoom("TEST")

// 		p1 := game.NewPlayer("p1", nil)
// 		p2 := game.NewPlayer("p2", nil)
// 		room.PlayerRegister(p1)
// 		room.PlayerRegister(p2)

// 		room.RemovePlayer(p1)

// 		if room.IsEmpty() {
// 			t.Error("Expected room one player after disconnect to be not empty")
// 		}
// 	})

// 	t.Run("room with two player after all player disconnected should be empty", func(t *testing.T) {
// 		room := game.NewGameRoom("TEST")

// 		p1 := game.NewPlayer("p1", nil)
// 		p2 := game.NewPlayer("p2", nil)
// 		room.PlayerRegister(p1)
// 		room.PlayerRegister(p2)

// 		room.RemovePlayer(p1)
// 		room.RemovePlayer(p2)
// 		if !room.IsEmpty() {
// 			t.Error("Expected room with two player after all player disconnected to be empty")
// 		}
// 	})
// }
