package game

import (
	"sync"
)

type GameState interface {
	StateInit()
	OnPlayerJoin(*Player)
	OnPlayerLeft(*Player)
	GetName() string
}

func GetStateName(room *GameRoom) string {
	return room.State.GetName()
}

type LobbyState struct {
	stateName    string
	playerStatus map[*Player]bool //True = Ready, False = Unready
	mu           sync.Mutex
}

func (l *LobbyState) StateInit() {
	l.stateName = "LOBBY_STATE"
	l.playerStatus = make(map[*Player]bool)
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

func (l *LobbyState) GetName() string {
	return "LOBBY_STATE"
}

func (l *LobbyState) PlayerReady(player *Player) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if exists := l.playerStatus[player]; exists {
		l.playerStatus[player] = true
	}
}

func (l *LobbyState) PlayerUnready(player *Player) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if exists := l.playerStatus[player]; exists {
		l.playerStatus[player] = false
	}
}
