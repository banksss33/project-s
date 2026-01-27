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
	GameClose      chan bool
	mu             sync.Mutex
	countdown      chan int
}

// create new game room Constructor
func NewGameRoom(gameClose chan bool, host *Player) *GameRoom {
	newRoom := &GameRoom{
		PlayerList:     make(map[string]*Player),
		state:          "LOBBY_STATE",
		Broadcast:      make(chan types.ServerResponse, 5),
		ActionReceiver: make(chan types.PlayerAction, 5),
		GameClose:      gameClose,
		countdown:      make(chan int),
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
	// Mock Data supposed to be from database or config
	var mockLocations = map[string][]string{
		"Airplane":      {"Pilot", "Co-pilot", "Flight Attendant", "Passenger", "Air Marshal", "Mechanic"},
		"Hospital":      {"Surgeon", "Nurse", "Patient", "Therapist", "Security Guard", "Doctor"},
		"Bank":          {"Bank Manager", "Teller", "Robber", "Customer", "Armored Car Driver", "Consultant"},
		"Pirate Ship":   {"Captain", "Cabin Boy", "Pirate", "Gunner", "Bound Prisoner", "Cook"},
		"Space Station": {"Astronaut", "Scientist", "Commander", "Engineer", "Alien", "Space Tourist"},
	}

	var lobby = &LobbyState{}
	lobby.setting.Locations = mockLocations
	var inGame = &InGameState{}

	lobbyAction := LobbyActionDispatcher()
	inGameAction := InGameActionDispatcher()
	timeUp := make(chan bool)
	for {
		select {
		case action, ok := <-gr.ActionReceiver:
			if !ok {
				break
			}

			if action.ActionName == "NEW_ROUND" {
				//Broadcast to all player about everyone roles before start new Round

				gr.state = "IN_GAME_STATE"
				inGame.NewRound()
				inGame.StartNewTimer(gr.countdown, timeUp)
				broadcastPlayerRolesAndLocation(inGame)
				broadcastGameState(gr, inGame)
			}

			if gr.state == "CLEARING_OLD_ACTION" {
				continue
			}

			if action.ActionName == "START_GAME" && gr.state == "LOBBY_STATE" && gr.PlayerList[action.UserID] == lobby.Host {
				gr.state = "IN_GAME_STATE"
				players := lobby.GetPlayers()
				spectators := lobby.GetSpectators()
				inGame.StartNewGame(lobby.setting, players, spectators)
				inGame.StartNewTimer(gr.countdown, timeUp)
				broadcastPlayerRolesAndLocation(inGame)
				broadcastGameState(gr, inGame)
			}

			switch gr.state {
			case "LOBBY_STATE":
				useLobbyAction := lobbyAction[action.ActionName]
				useLobbyAction(gr, lobby, action)
			case "IN_GAME_STATE":
				useInGameAction := inGameAction[action.ActionName]
				useInGameAction(gr, inGame, action)
				if inGame.isRoundEnd {
					gr.state = "CLEARING_OLD_ACTION"
					inGame.killTimer <- true
					gr.ActionReceiver <- types.PlayerAction{
						ActionName: "NEW_ROUND",
					}
				}
			}
		case cd := <-gr.countdown:
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
		case <-timeUp:
			gr.state = "CLEARING_OLD_ACTION"
			gr.ActionReceiver <- types.PlayerAction{
				ActionName: "NEW_ROUND",
			}
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
