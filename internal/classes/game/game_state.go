package game

import (
	"math/rand/v2"
	"project-s/internal/types"
	"sync"
	"time"
)

type LobbyState struct {
	stateName   string
	Host        *Player          //player that are host
	readyStatus map[*Player]bool //True = Ready, False = Unready
	mu          sync.Mutex

	types.GameSetting
}

func (l *LobbyState) Init(spiesCount, timer int, locations map[string][]string) {
	l.stateName = "LOBBY_STATE"
	l.readyStatus = make(map[*Player]bool)

	//init default setting
	l.Spies = spiesCount
	l.Timer = timer
	l.Locations = locations
}

func (l *LobbyState) OnPlayerJoin(player *Player) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.readyStatus[player] = false
}

func (l *LobbyState) OnPlayerLeft(player *Player) {
	l.mu.Lock()
	defer l.mu.Unlock()

	delete(l.readyStatus, player)
}

func (l *LobbyState) Name() string {
	return l.stateName
}

func (l *LobbyState) PlayerReady(player *Player) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, exists := l.readyStatus[player]; exists {
		l.readyStatus[player] = true
	}
}

func (l *LobbyState) PlayerUnready(player *Player) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, exists := l.readyStatus[player]; exists {
		l.readyStatus[player] = false
	}
}

func (l *LobbyState) EditSetting(newSetting types.GameSetting) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.Locations = newSetting.Locations
	l.Spies = newSetting.Spies
	l.Timer = newSetting.Timer
}

func (l *LobbyState) IsReady() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if len(l.readyStatus) < 3 {
		return false
	}

	for _, ready := range l.readyStatus {
		if !ready {
			return false
		}
	}

	return true
}

type InGameState struct {
	timer       types.GameTimer
	stateName   string
	location    string
	isVoting    bool
	playerRoles map[*Player]types.PlayerStatus //key: player | value: player's roles
	mu          sync.Mutex
}

// random player roles here
func (i *InGameState) Init(setting types.GameSetting, playerList map[*Player]bool) {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.stateName = "IN_GAME_STATE"

	i.timer = types.GameTimer{
		Tick:      time.NewTicker(time.Second),
		IsRunning: false,
		Countdown: setting.Timer,
	}

	i.isVoting = false
	i.playerRoles = make(map[*Player]types.PlayerStatus)

	//random location
	randomLocation := randomkeyFromMap[string, []string](setting.Locations, 1)[0]

	//random roles to player
	shuffledRoles := setting.Locations[randomLocation]
	rand.Shuffle(len(randomLocation), func(a, b int) {
		shuffledRoles[a], shuffledRoles[b] = shuffledRoles[b], shuffledRoles[a]
	})

	//assign random spies role from playerList then remove that from playerList
	randomSpies := randomkeyFromMap(playerList, uint(setting.Spies))
	for _, player := range randomSpies {
		i.playerRoles[player] = types.PlayerStatus{
			Score: 0,
			Roles: "SPY",
		}
		delete(playerList, player)
	}

	//assign random location to player
	index := 0
	for player, _ := range playerList {
		i.playerRoles[player] = types.PlayerStatus{
			Score: 0,
			Roles: shuffledRoles[index],
		}
		index++
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
	//

}

func (i *InGameState) OnPlayerLeft(player *Player) {

}

func (i *InGameState) Name() string {
	return i.stateName
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
