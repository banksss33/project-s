package game

import (
	"math/rand/v2"
	"project-s/internal/types"
	"slices"
	"time"
)

type InGameState struct {
	roundLeft int
	timer     types.GameTimer
	location  string

	killTimer     chan bool
	isVoting      bool
	isRoundEnd    bool
	voteLeft      int
	accusedPlayer *Player

	playerStatus        map[*Player]*types.PlayerStatus //key: player | value: player's roles
	spectator           map[*Player]bool
	disconnectedPlayers map[*Player]bool
	setting             types.GameSetting
}

func NewInGameState() *InGameState {
	return &InGameState{
		spectator:           make(map[*Player]bool),
		killTimer:           make(chan bool),
		playerStatus:        make(map[*Player]*types.PlayerStatus),
		disconnectedPlayers: make(map[*Player]bool),
		isVoting:            false,
		isRoundEnd:          false,
		accusedPlayer:       nil,
		timer: types.GameTimer{
			Tick:      time.NewTicker(time.Second),
			IsRunning: false,
			Countdown: 0,
		},
	}
}

func (i *InGameState) StartNewGame(gameSetting types.GameSetting, playerList []*Player, spectatorList []*Player) {
	i.setting = gameSetting
	for _, s := range spectatorList {
		i.spectator[s] = true
	}
	i.roundLeft = i.setting.Round

	i.timer.Countdown = i.setting.Timer

	//random location
	randomLocation := randomkeyFromMap[string, []string](i.setting.Locations, 1)[0]
	i.location = randomLocation

	//random roles to player
	shuffledPlayer := slices.Clone(playerList)
	rand.Shuffle(len(shuffledPlayer), func(a, b int) {
		shuffledPlayer[a], shuffledPlayer[b] = shuffledPlayer[b], shuffledPlayer[a]
	})

	roles := i.setting.Locations[randomLocation]
	for index, player := range shuffledPlayer {
		if index < i.setting.Spies {
			i.playerStatus[player] = &types.PlayerStatus{
				Score:       0,
				Roles:       "SPY",
				AlreadyVote: false,
			}

			continue
		}

		i.playerStatus[player] = &types.PlayerStatus{
			Score:       0,
			Roles:       roles[index-i.setting.Spies],
			AlreadyVote: false,
		}
	}
}

func (i *InGameState) NewRound() {
	i.timer.Countdown = i.setting.Timer
	i.isVoting = false
	i.isRoundEnd = false
	i.accusedPlayer = nil

	var inGamePlayer []*Player
	for player := range i.playerStatus {
		inGamePlayer = append(inGamePlayer, player)
	}

	randomLocation := randomkeyFromMap[string, []string](i.setting.Locations, 1)[0]
	i.location = randomLocation

	//random roles to player
	shuffledPlayer := slices.Clone(inGamePlayer)
	rand.Shuffle(len(shuffledPlayer), func(a, b int) {
		shuffledPlayer[a], shuffledPlayer[b] = shuffledPlayer[b], shuffledPlayer[a]
	})

	roles := i.setting.Locations[randomLocation]
	for index, player := range shuffledPlayer {

		if index < i.setting.Spies {
			i.playerStatus[player].Roles = "SPY"
			i.playerStatus[player].AlreadyVote = false

			continue
		}

		i.playerStatus[player].Roles = roles[(index-i.setting.Spies)%len(roles)]
		i.playerStatus[player].AlreadyVote = false

	}
}

func (i *InGameState) StartNewTimer(countdown chan<- int, timeUp chan<- bool) {
	i.timer.Tick.Reset(time.Second)
	i.timer.IsRunning = true

	go func() {
		defer i.timer.Tick.Stop()

		for {
			select {
			case <-i.killTimer:
				return
			case <-i.timer.Tick.C:
				if i.timer.Countdown < 1 {

					for _, status := range i.playerStatus {
						if status.Roles == "SPY" {
							status.Score += 2
						}
					}
					i.roundLeft--
					i.isRoundEnd = true
					timeUp <- true
					return
				}
				countdown <- i.timer.Countdown
				i.timer.Countdown--
			}
		}
	}()
}

func (i *InGameState) ResumeTimer() {
	if i.timer.IsRunning {
		return
	}

	i.timer.Tick.Reset(time.Second)
	i.timer.IsRunning = true
}

func (i *InGameState) PauseTimer() {
	if !i.timer.IsRunning {
		return
	}

	i.timer.Tick.Stop()
	i.timer.IsRunning = false
}

func (i *InGameState) OnPlayerJoin(player *Player) {
	// If player is not in playerStatus, they are a new join -> spectator
	if _, exists := i.disconnectedPlayers[player]; !exists {
		i.spectator[player] = true
	}

	delete(i.disconnectedPlayers, player)

	// Resume timer if no more disconnected players
	if len(i.disconnectedPlayers) < 1 {
		i.ResumeTimer()
	}
}

func (i *InGameState) OnPlayerLeave(player *Player) (shouldCloseRoom bool) {
	if _, exists := i.spectator[player]; exists {
		delete(i.spectator, player)
		return false
	}

	if _, exists := i.playerStatus[player]; exists {
		i.disconnectedPlayers[player] = true
		i.PauseTimer()

		if i.AllPlayersDisconnected() {
			return true
		}
	}

	return false
}

func (i *InGameState) AllPlayersDisconnected() bool {
	if len(i.disconnectedPlayers) == len(i.playerStatus) {
		return true
	}

	return false
}

func (i *InGameState) Accuse(fromPlayer *Player, targetPlayer *Player) {
	_, accruserExists := i.playerStatus[fromPlayer]
	_, accrusedExists := i.playerStatus[targetPlayer]
	if !accruserExists || !accrusedExists {
		return
	}

	if i.accusedPlayer == nil || fromPlayer == i.accusedPlayer {
		return
	}

	i.isVoting = true
	i.voteLeft = len(i.playerStatus) - 1
	i.accusedPlayer = targetPlayer

	i.Vote(fromPlayer)
}

func (i *InGameState) Vote(fromPlayer *Player) {
	voter, voterExists := i.playerStatus[fromPlayer]
	if !i.isVoting || i.accusedPlayer == nil || fromPlayer == i.accusedPlayer || voter.AlreadyVote || !voterExists {
		return
	}

	voter.AlreadyVote = true
	i.voteLeft--

	if i.voteLeft > 0 {
		return
	}

	//vote success Round end
	accrusedPlayerStatus := i.playerStatus[i.accusedPlayer]
	if accrusedPlayerStatus.Roles == "SPY" {
		for _, status := range i.playerStatus {
			if status.Roles != "SPY" {
				status.Score += 2
			}
		}
	} else {
		for _, status := range i.playerStatus {
			if status.Roles == "SPY" {
				status.Score += 4
			}
		}
	}

	i.roundLeft--
	i.isRoundEnd = true
}

func (i *InGameState) Disagree(fromPlayer *Player) {
	playerDisagree, playerDisagreeExists := i.playerStatus[fromPlayer]
	if !i.isVoting || i.accusedPlayer == nil || fromPlayer == i.accusedPlayer || playerDisagree.AlreadyVote || !playerDisagreeExists {
		return
	}

	i.isVoting = false
	i.accusedPlayer = nil
	i.voteLeft = 0
	for _, playerStatus := range i.playerStatus {
		playerStatus.AlreadyVote = false
	}
}

func (i *InGameState) SpyVoteLocation(fromPlayer *Player, Location string) {
	if status := i.playerStatus[fromPlayer]; status.Roles != "SPY" {
		return
	}

	if Location == i.location {
		for _, status := range i.playerStatus {
			if status.Roles == "SPY" {
				status.Score += 4
			}
		}
	} else {
		for _, status := range i.playerStatus {
			if status.Roles != "SPY" {
				status.Score += 1
			}
		}
	}

	i.roundLeft--
	i.isRoundEnd = true
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
