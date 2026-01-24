package game

import (
	"project-s/internal/types"
	"sync"
)

type LobbyState struct {
	host       *Player          //player that are host
	playerList map[*Player]bool //Key:  player, value: true = player, false = spectator
	mu         sync.Mutex

	setting types.GameSetting
}

func (l *LobbyState) Init(host *Player, locations map[string][]string) {
	l.playerList = make(map[*Player]bool)

	//init default setting
	l.setting.Spies = 1
	l.setting.Timer = 420
	l.setting.Locations = locations
	l.host = host
}

func (l *LobbyState) OnPlayerJoin(player *Player) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.playerList[player] = true
}

func (l *LobbyState) OnSpectatorJoin(player *Player) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.playerList[player] = false
}

func (l *LobbyState) OnPlayerLeft(player *Player) {
	l.mu.Lock()
	defer l.mu.Unlock()

	delete(l.playerList, player)
}

func (l *LobbyState) EditSetting(newSetting types.GameSetting) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.Locations = newSetting.Locations
	l.Spies = newSetting.Spies
	l.Timer = newSetting.Timer
}
