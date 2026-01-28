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
	ActionReceiver chan types.Action
	GameClose      chan bool
	mu             sync.Mutex
	countdown      chan int
	timeUp         chan bool
}

// create new game room Constructor
func NewGameRoom(gameClose chan bool, host *Player) *GameRoom {
	newRoom := &GameRoom{
		PlayerList:     make(map[string]*Player),
		state:          "LOBBY_STATE",
		Broadcast:      make(chan types.ServerResponse, 5),
		ActionReceiver: make(chan types.Action, 5),
		GameClose:      gameClose,
		countdown:      make(chan int),
		timeUp:         make(chan bool),
	}

	go newRoom.broadcastInit()
	go newRoom.actionProcessorInit()

	return newRoom
}

func (gr *GameRoom) broadcastInit() {
	for broadcastItem := range gr.Broadcast {
		for _, player := range gr.PlayerList {
			if player.IsConnected {
				select {
				case player.SendToPlayer <- broadcastItem:
				default:
				}
			}
		}
	}
}

func (gr *GameRoom) actionProcessorInit() {
	var MockLocations = map[string][]string{
		"Airplane":      {"Pilot", "Co-pilot", "Flight Attendant", "Passenger", "Air Marshal", "Mechanic"},
		"Hospital":      {"Surgeon", "Nurse", "Patient", "Therapist", "Security Guard", "Doctor"},
		"Bank":          {"Bank Manager", "Teller", "Robber", "Customer", "Armored Car Driver", "Consultant"},
		"Pirate Ship":   {"Captain", "Cabin Boy", "Pirate", "Gunner", "Bound Prisoner", "Cook"},
		"Space Station": {"Astronaut", "Scientist", "Commander", "Engineer", "Alien", "Space Tourist"},
	}
	var lobby *LobbyState = NewLobbyState()
	var inGame *InGameState

	lobby.setting.Locations = MockLocations

	lobbyAction := LobbyActionDispatcher()
	inGameAction := InGameActionDispatcher()
	transitionAction := GameTransitionActionDispatcher()

	for {
		select {
		case action, ok := <-gr.ActionReceiver:
			if !ok {
				break
			}

			switch gr.state {
			case "LOBBY_STATE":
				if action.ActionName == "START_GAME" {
					player, _ := gr.GetPlayerByID(action.CallerID)
					if lobby.Host != nil && player == lobby.Host {
						gr.state = "IN_GAME_STATE"
						players := lobby.GetPlayers()
						spectators := lobby.GetSpectators()
						inGame = NewInGameState()
						inGame.StartNewGame(lobby.setting, players, spectators)
						inGame.StartNewTimer(gr.countdown, gr.timeUp)
						SendPlayerRolesAndLocation(inGame)
						broadcastGameState(gr, inGame)
					}
					continue
				}

				useLobbyAction, lobbyActionExists := lobbyAction[action.ActionName]
				if lobbyActionExists {
					useLobbyAction(gr, lobby, action)
				}

			case "IN_GAME_STATE":
				useInGameAction, useInGameActionExists := inGameAction[action.ActionName]
				if useInGameActionExists {
					useInGameAction(gr, inGame, action)
				}

				if inGame.isRoundEnd {
					gr.state = "TRANSITION_STATE"
					inGame.killTimer <- true

					serverAction := "NEW_ROUND"
					if inGame.roundLeft < 1 {
						serverAction = "BACK_TO_LOBBY"
					}

					gr.ActionReceiver <- types.Action{
						CallerID:   "SERVER",
						ActionName: serverAction,
					}
				}

			case "TRANSITION_STATE":
				useTransitionAction, useTransitionActionExists := transitionAction[action.ActionName]
				if useTransitionActionExists {
					useTransitionAction(gr, inGame, lobby, action)
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
		case <-gr.timeUp:
			gr.state = "TRANSITION_STATE"

		}
	}
}

func (gr *GameRoom) GetPlayerByID(userID string) (*Player, bool) {
	gr.mu.Lock()
	defer gr.mu.Unlock()
	player, exists := gr.PlayerList[userID]

	return player, exists
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
	gr.mu.Lock()
	defer gr.mu.Unlock()
	for _, player := range gr.PlayerList {
		if player.IsConnected {
			return false
		}
	}
	return true
}
