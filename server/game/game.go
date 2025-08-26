package game

import "github.com/gofiber/contrib/websocket"

const StatusConnected = "CONNECTED"
const StatusDisconnected = "DISCONNECTED"

// #region Player
type Player struct {
	ConnectionStatus string // CONNECTED | DISCONNECTED
	Conn             *websocket.Conn
	Score            uint // Score
}

// #region GameSetting
type GameSetting struct {
	RoundCount uint // round count before the game end | default as 7
	Timer      uint // seconds | default as 420 seconds -> 7 min
	Spies      uint // can only be 1 or 2 | default as 1
}

// #region GameRoom
type GameRoom struct {
	PlayerList map[string]*Player // Key: userID from Discord, Value: *Player
	Setting    GameSetting
	State      string
}

// create new game room
func NewGameRoom() *GameRoom {
	return &GameRoom{
		PlayerList: make(map[string]*Player),
		Setting:    GameSetting{6, 420, 1},
		State:      "",
	}
}

//add new player to game room
func (gr *GameRoom) AddPlayer(userID string, conn *websocket.Conn) {
	gr.PlayerList[userID] = &Player{
		ConnectionStatus: StatusConnected,
		Conn:             conn,
		Score:            0,
	}
}

//use when player reconnected
func (gr *GameRoom) PlayerReconnected(userID string, conn *websocket.Conn) {
	reconnectedPlayer := gr.PlayerList[userID]
	reconnectedPlayer.ConnectionStatus = StatusConnected
	reconnectedPlayer.Conn = conn //assign new connection
}

//use when player disconnected
func (gr *GameRoom) PlayerDisconnected(userID string) {
	disconnectedPlayer := gr.PlayerList[userID]
	disconnectedPlayer.ConnectionStatus = StatusDisconnected
	disconnectedPlayer.Conn = nil
}

//If all player disconnected return true
func (gr *GameRoom) IsRoomEmpty() bool {
	for _, player := range gr.PlayerList {
		if player.ConnectionStatus == StatusConnected {
			return false
		}
	}
	return true
}
