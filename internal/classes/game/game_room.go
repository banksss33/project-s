package game

import (
	"encoding/json"
	"project-s/internal/types"
	"sync"
)

type GameRoom struct {
	state          string
	PlayerList     map[string]*Player // Key: userID from Discord, Value: *Player
	Broadcast      chan types.ServerResponse
	ActionReceiver chan types.PlayerAction
	gameClose      chan bool
	mu             sync.Mutex
}

// create new game room Constructor
func NewGameRoom(gameClose chan bool, host *Player) *GameRoom {
	newRoom := &GameRoom{
		PlayerList:     make(map[string]*Player),
		state:          "LOBBY_STATE",
		Broadcast:      make(chan types.ServerResponse, 5),
		ActionReceiver: make(chan types.PlayerAction, 5),
		gameClose:      gameClose,
	}

	go newRoom.broadcastInit()
	go newRoom.actionProcessorInit()

	return newRoom
}

func (gr *GameRoom) broadcastInit() {
	for broadcastItem := range gr.Broadcast {
		for _, player := range gr.PlayerList {
			player.SendToPlayer <- broadcastItem
		}
	}
}

func (gr *GameRoom) actionProcessorInit() {
	countdown := make(chan int)
	var lobby = &LobbyState{}
	var inGame = &InGameState{}

	lobbyAction := LobbyActionDispatcher()
	inGameAction := InGameActionDispatcher()
	for {
		select {
		case action, ok := <-gr.ActionReceiver:
			if !ok {
				break
			}

			if action.ActionName == "START_GAME" && gr.state == "LOBBY_STATE" && gr.PlayerList[action.UserID] == lobby.Host {
				gr.state = "IN_GAME_STATE"
				players := lobby.GetPlayers()
				spectators := lobby.GetSpectators()
				inGame.Init(lobby.setting, players, spectators)
				inGame.StartTimer(countdown)
			}

			switch gr.state {
			case "LOBBY_STATE":
				useLobbyAction := lobbyAction[action.ActionName]
				useLobbyAction(gr, lobby, action)
			case "IN_GAME_STATE":
				useInGameAction := inGameAction[action.ActionName]
				useInGameAction(gr, inGame, action)
			}
		case cd := <-countdown:
			//broadcast time
			timerUpdate := types.UpdateTimerResponse{
				TimerNow: cd,
			}

			jsonTimerUpdatePayload, _ := json.Marshal(timerUpdate)
			serverResponse := types.ServerResponse{
				ResponseName: "TIMER_UPDATE",
				Payload:      jsonTimerUpdatePayload,
			}

			gr.Broadcast <- serverResponse
		}
	}
}

func (gr *GameRoom) GetPlayerByID(userID string) *Player {
	gr.mu.Lock()
	defer gr.mu.Unlock()
	player := gr.PlayerList[userID]

	return player
}

// add new player to game room
func (gr *GameRoom) PlayerRegister(player *Player) {
	gr.mu.Lock()
	defer gr.mu.Unlock()

	gr.PlayerList[player.UserID] = player
}

func (gr *GameRoom) Cleanup() {
	if !gr.isEmpty() {
		return
	}

	close(gr.ActionReceiver)
	close(gr.Broadcast)
}

//#region action method - helper

// If all player disconnected return true
func (gr *GameRoom) isEmpty() bool {
	for _, player := range gr.PlayerList {
		if player.IsConnected {
			return false
		}
	}
	return true
}
