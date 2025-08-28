package game_test

import (
	"project-s/internal/classes/game"
	"testing"
)

func TestIsRoomEmpty(t *testing.T) {
	t.Run("newly created room should be empty", func(t *testing.T) {
		room := game.NewGameRoom()

		if !room.IsRoomEmpty() {
			t.Error("Expected new room to be empty")
		}
	})

	t.Run("room with one connected player should not be empty", func(t *testing.T) {
		room := game.NewGameRoom()
		room.AddPlayer("p1", nil)

		if room.IsRoomEmpty() {
			t.Error("Expected room with one player to not be empty")
		}
	})

	t.Run("room after all players disconnect should be empty", func(t *testing.T) {
		room := game.NewGameRoom()
		room.AddPlayer("p1", nil)
		room.PlayerDisconnected("p1")

		if !room.IsRoomEmpty() {
			t.Error("Expected room to be empty after player disconnected")
		}
	})

	t.Run("room with two player should not be empty", func(t *testing.T) {
		room := game.NewGameRoom()

		room.AddPlayer("p1", nil)
		room.AddPlayer("p2", nil)

		if room.IsRoomEmpty() {
			t.Error("Expected room with two players to not be empty")
		}
	})

	t.Run("room with one player after player disconnect should not be empty", func(t *testing.T) {
		room := game.NewGameRoom()

		room.AddPlayer("p1", nil)
		room.AddPlayer("p2", nil)

		room.PlayerDisconnected("p1")

		if room.IsRoomEmpty() {
			t.Error("Expected room one player after disconnect to be not empty")
		}
	})

	t.Run("room with two player after all player disconnected should be empty", func(t *testing.T) {
		room := game.NewGameRoom()
		room.AddPlayer("p1", nil)
		room.AddPlayer("p2", nil)

		room.PlayerDisconnected("p1")
		room.PlayerDisconnected("p2")
		if !room.IsRoomEmpty() {
			t.Error("Expected room with two player after all player disconnected to be empty")
		}
	})
}

func TestNewGameRoom(t *testing.T) {
	t.Run("new game room setting default round should be 6 round", func(t *testing.T) {
		room := game.NewGameRoom()
		setting := room.Setting
		if setting.RoundCount != 6 {
			t.Error("Expect default round count to be 6")
		}
	})

	t.Run("new game room setting default timer should be 420 seconds", func(t *testing.T) {
		room := game.NewGameRoom()
		setting := room.Setting
		if setting.Timer != 420 {
			t.Error("Expect default timer to be 420")
		}
	})

	t.Run("new game room setting default spies should be 1", func(t *testing.T) {
		room := game.NewGameRoom()
		setting := room.Setting
		if setting.Spies != 1 {
			t.Error("Expect default spies to be 1")
		}
	})
}
