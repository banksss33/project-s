package game

import (
	"encoding/json"
)

type GameRoom struct {
	RoomID     string
	PlayerList map[string]*Player // Key: userID from Discord, Value: *Player
	State      GameState
	Broadcast  chan json.RawMessage
}

// create new game room Constructor
func NewGameRoom(roomID string) *GameRoom {
	newRoom := &GameRoom{
		RoomID:     roomID,
		PlayerList: make(map[string]*Player),
		State:      &LobbyState{},
		Broadcast:  make(chan json.RawMessage),
	}

	newRoom.State.StateInit()
	go newRoom.BroadcastInit()

	return newRoom
}

func (gr *GameRoom) BroadcastInit() {
	for broadcastItem := range gr.Broadcast {
		for _, player := range gr.PlayerList {
			player.SendJSON <- broadcastItem
		}
	}
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
