package game

import (
	"project-s/internal/types"
)

type LobbyState struct {
	Host       *Player          //player that are host
	PlayerList map[*Player]bool //Key:  player, value: true = player, false = spectator
	setting    types.GameSetting
}

func NewLobbyState() *LobbyState {
	return &LobbyState{
		PlayerList: make(map[*Player]bool),
		setting: types.GameSetting{
			Round:     7,
			Spies:     1,
			Timer:     420,
			Locations: make(map[string][]string),
		},
	}
}

func (l *LobbyState) PlayerJoin(player *Player) {
	l.PlayerList[player] = true
	if l.Host == nil {
		l.Host = player
	}
}

func (l *LobbyState) SpectatorJoin(player *Player) {
	if !l.PlayerList[player] { //case when player are already spectator
		return
	}

	if player == l.Host {
		for randPlayer, isPlayer := range l.PlayerList {
			if !isPlayer {
				continue
			}

			if randPlayer != l.Host {
				l.Host = randPlayer
				break
			}
		}
	}

	l.PlayerList[player] = false

	//if all player are spectator then Host = nil
	for _, isPlayer := range l.PlayerList {
		// if found player in lobby function end
		if isPlayer {
			return
		}
	}

	//if all player are spectator then host is nil
	l.Host = nil
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
	var players []*Player
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
