package game

import (
	"project-s/internal/types"
)

type LobbyState struct {
	host       *Player          //player that are host
	playerList map[*Player]bool //Key:  player, value: true = player, false = spectator
	setting    types.GameSetting
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
	l.playerList[player] = true
}

func (l *LobbyState) OnSpectatorJoin(player *Player) {
	l.playerList[player] = false
}

func (l *LobbyState) OnPlayerLeft(player *Player) {
	delete(l.playerList, player)
}

func (l *LobbyState) EditSetting(newSetting types.GameSetting) {
	l.setting.Locations = newSetting.Locations
	l.setting.Spies = newSetting.Spies
	l.setting.Timer = newSetting.Timer
}
