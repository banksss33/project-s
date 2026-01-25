package game

import (
	"project-s/internal/types"
)

type LobbyState struct {
	Host       *Player          //player that are host
	PlayerList map[*Player]bool //Key:  player, value: true = player, false = spectator
	setting    types.GameSetting
}

func (i *LobbyState) LobbyStateInit(host *Player, locations map[string][]string) {
	i.Host = host
	i.PlayerList = make(map[*Player]bool)
	i.setting = types.GameSetting{
		Round:     5,
		Spies:     1,
		Timer:     420,
		Locations: locations,
	}

}

func (l *LobbyState) PlayerJoin(player *Player) {
	l.PlayerList[player] = true
}

func (l *LobbyState) SpectatorJoin(player *Player) {
	if player == l.Host {
		for randPlayer := range l.PlayerList {
			if randPlayer != player {
				l.Host = randPlayer
				break
			}
		}
	}

	l.PlayerList[player] = false
}

func (l *LobbyState) PlayerLeft(player *Player) {
	if player == l.Host {
		for randPlayer := range l.PlayerList {
			if randPlayer != player {
				l.Host = randPlayer
				break
			}
		}
	}

	delete(l.PlayerList, player)
}

func (l *LobbyState) EditSetting(newSetting types.GameSetting) {
	l.setting.Locations = newSetting.Locations
	l.setting.Spies = newSetting.Spies
	l.setting.Timer = newSetting.Timer
}
