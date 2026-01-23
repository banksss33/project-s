package game

import (
	"math/rand/v2"
	"project-s/internal/types"
	"slices"
	"sync"
	"time"
)

type InGameState struct {
	roundLeft    int
	timer        types.GameTimer
	location     string
	isVoting     bool
	playerStatus map[*Player]types.PlayerStatus //key: player | value: player's roles
	spectator    []*Player
	mu           sync.Mutex
}

// random player roles here
func (i *InGameState) Init(setting types.GameSetting, playerList []*Player, spectatorList []*Player) {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.roundLeft = setting.Round
	i.spectator = slices.Clone(spectatorList)
	i.timer = types.GameTimer{
		Tick:      time.NewTicker(time.Second),
		IsRunning: false,
		Countdown: setting.Timer,
	}
	i.isVoting = false
	i.playerStatus = make(map[*Player]types.PlayerStatus)

	//random location
	randomLocation := randomkeyFromMap[string, []string](setting.Locations, 1)[0]
	i.location = randomLocation

	//random roles to player

	shuffledPlayer := slices.Clone(playerList)
	rand.Shuffle(len(shuffledPlayer), func(a, b int) {
		shuffledPlayer[a], shuffledPlayer[b] = shuffledPlayer[b], shuffledPlayer[a]
	})

	roles := setting.Locations[randomLocation]
	for index, player := range playerList {
		if index < setting.Spies {
			i.playerStatus[player] = types.PlayerStatus{
				Score: 0,
				Roles: "SPY",
			}

			continue
		}

		i.playerStatus[player] = types.PlayerStatus{
			Score: 0,
			Roles: roles[index-setting.Spies],
		}
	}

}

func (i *InGameState) StartTimer(countdown chan<- int) {
	if i.timer.IsRunning {
		return
	}

	i.timer.IsRunning = true
	go func() {
		for i.timer.Countdown > 0 {
			<-i.timer.Tick.C
			i.timer.Countdown--
			countdown <- i.timer.Countdown
		}
	}()
}

func (i *InGameState) PauseTimer() {
	if !i.timer.IsRunning {
		return
	}

	i.timer.Tick.Stop()
	i.timer.IsRunning = false
}

func (i *InGameState) OnPlayerJoin(player *Player) {
	//check if player exists in playerRoles if not then it mean recently join = spectator
	if _, exists := i.playerStatus[player]; exists {
		i.playerStatus[player] = types.PlayerStatus{
			Score: 0,
			Roles: "SPECTATOR",
		}
	}
}

func (i *InGameState) Vote(fromPlayer *Player, targetPlayer *Player) {
	if i.isVoting {
		return
	}

	i.isVoting = true

}

// helper func
func randomkeyFromMap[K comparable, V any](paramMap map[K]V, retCount uint) []K {
	paramLength := len(paramMap)
	randomkey := make([]K, 0, paramLength)
	for key, _ := range paramMap {
		randomkey = append(randomkey, key)
	}

	rand.Shuffle(paramLength, func(a, b int) {
		randomkey[a], randomkey[b] = randomkey[b], randomkey[a]
	})

	return randomkey[0:retCount]
}
