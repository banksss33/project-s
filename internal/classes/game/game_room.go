package game

import (
	"project-s/internal/services/pubsub"
)

type GameRoom struct {
	RoomID     string
	PlayerList map[string]*Player // Key: userID from Discord, Value: *Player
	Publish    *pubsub.Publisher
	State      GameState
}

// create new game room Constructor
func NewGameRoom(roomID string) *GameRoom {
	newRoom := &GameRoom{
		RoomID:     roomID,
		PlayerList: make(map[string]*Player),
		Publish:    pubsub.NewPublisher(),
		State:      &LobbyState{},
	}

	newRoom.State.StateInit()

	return newRoom
}

// add new player to game room
func (gr *GameRoom) AddPlayer(player *Player) {
	defer gr.State.OnPlayerJoin(player)
	if !player.IsConnected { //case for new player
		gr.PlayerList[player.UserID] = player
		return
	}

}

// If all player disconnected return true
func (gr *GameRoom) IsRoomEmpty() bool {
	for _, player := range gr.PlayerList {
		if !player.IsConnected {
			return false
		}
	}
	return true
}
