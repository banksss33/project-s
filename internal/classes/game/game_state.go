package game

import (
	"sync"
)

type GameState interface {
	OnPlayerJoin(*Player)
	OnPlayerLeft(*Player)
	Name() string
}

func GetStateName(room *GameRoom) string {
	return room.State.Name()
}

type LobbyState struct {
	stateName    string
	playerStatus map[*Player]bool //True = Ready, False = Unready
	mu           sync.Mutex

	//setting
	spies     int
	timer     int
	locations map[string][]string //key: location name| value: location roles
}

func (l *LobbyState) Init() {
	l.stateName = "LOBBY_STATE"
	l.playerStatus = make(map[*Player]bool)

	//init default setting
	l.spies = 1
	l.timer = 420
	l.locations = make(map[string][]string)
}

func (l *LobbyState) OnPlayerJoin(player *Player) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.playerStatus[player] = false

}

func (l *LobbyState) OnPlayerLeft(player *Player) {
	l.mu.Lock()
	defer l.mu.Unlock()

	delete(l.playerStatus, player)
}

func (l *LobbyState) Name() string {
	return l.stateName
}

func (l *LobbyState) PlayerReady(player *Player) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, exists := l.playerStatus[player]; exists {
		l.playerStatus[player] = true
	}
}

func (l *LobbyState) PlayerUnready(player *Player) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, exists := l.playerStatus[player]; exists {
		l.playerStatus[player] = false
	}
}

func (l *LobbyState) addLocation(locationName string, locationRoles []string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, exists := l.locations[locationName]; exists {
		return
	}

	l.locations[locationName] = locationRoles
}

func (l *LobbyState) EditLocation(targetLocation string, newRoles []string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, exists := l.locations[targetLocation]; exists {
		l.locations[targetLocation] = newRoles
	}
}

func (l *LobbyState) isReady() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if len(l.playerStatus) < 3 {
		return false
	}

	for _, ready := range l.playerStatus {
		if !ready {
			return false
		}
	}

	return true
}

type InGameState struct {
	stateName      string
	timerCountdown int
	location       string
	playerRoles    map[*Player]string //key: player | value: player's roles
}

func (i *InGameState) Init(timer int, location string, playerList map[*Player]bool) {
	i.stateName = "IN_GAME_STATE"
}

func (i *InGameState) OnPlayerJoin(player *Player) {

}

func (i *InGameState) OnPlayerLeft(player *Player) {

}

func (i *InGameState) Name() string {
	return i.stateName
}
