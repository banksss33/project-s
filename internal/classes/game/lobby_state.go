package game

import (
	"project-s/internal/types"
)

type LobbyState struct {
	Host       *Player          //player that are host
	PlayerList map[*Player]bool //Key:  player, value: true = player, false = spectator
	setting    types.GameSetting
}

func (i *LobbyState) LobbyStateInit(host *Player) {
	i.Host = host
	i.PlayerList = make(map[*Player]bool)
	i.setting.Round = 7
	i.setting.Spies = 1
	i.setting.Timer = 420
}

func (l *LobbyState) PlayerJoin(player *Player) {
	l.PlayerList[player] = true
	if l.Host == nil {
		l.Host = player
	}
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
	for _, isPlayer := range l.PlayerList {
		if isPlayer {
			break
		}

		l.Host = nil
	}
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

func (l *LobbyState) GetPlayers() []*Player {
	players := make([]*Player, 0)
	for player, isPlayer := range l.PlayerList {
		if isPlayer {
			players = append(players, player)
		}
	}

	return players
}
func (l *LobbyState) GetSpectators() []*Player {
	spectators := make([]*Player, 0)
	for player, isPlayer := range l.PlayerList {
		if !isPlayer {
			spectators = append(spectators, player)
		}
	}

	return spectators
}
