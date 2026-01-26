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

	isVoting      bool
	isRoundEnd    bool
	voteLeft      int
	accusedPlayer *Player

	playerStatus map[*Player]*types.PlayerStatus //key: player | value: player's roles
	spectator    []*Player
	setting      types.GameSetting
}

// random player roles here
func (i *InGameState) Init(gameSetting types.GameSetting, playerList []*Player, spectatorList []*Player) {
	i.setting = gameSetting
	i.spectator = slices.Clone(spectatorList)

	//move to start timer
	i.timer = types.GameTimer{
		Tick:      time.NewTicker(time.Second),
		IsRunning: false,
		Countdown: i.setting.Timer,
	}

	i.isVoting = false
	i.isRoundEnd = false
	i.accusedPlayer = nil
	i.playerStatus = make(map[*Player]*types.PlayerStatus)

	//below here can move to new function
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
	i.timer = types.GameTimer{
		Tick:      time.NewTicker(time.Second),
		IsRunning: false,
		Countdown: i.setting.Timer,
	}

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

		i.playerStatus[player].Roles = roles[index-i.setting.Spies]
		i.playerStatus[player].AlreadyVote = false

	}
}

func (i *InGameState) StartTimer(countdown chan<- int) {
	i.timer.IsRunning = true
	go func() {
		defer close(countdown)

		for i.timer.Countdown > 0 {
			<-i.timer.Tick.C
			countdown <- i.timer.Countdown
			i.timer.Countdown--
		}

		for _, status := range i.playerStatus {
			if status.Roles == "SPY" {
				status.Score += 2
			}
		}

		i.roundLeft--
		i.isRoundEnd = true
	}()
}

func (i *InGameState) ResumeTimer() {
	if i.timer.IsRunning {
		return
	}

	i.timer.Tick.Reset(time.Second)
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
		i.spectator = append(i.spectator, player)
	}
}

func (i *InGameState) Accruse(fromPlayer *Player, targetPlayer *Player) {
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

func (i *InGameState) SpyVoteLocation(Location string) {
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

func (i *InGameState) GetCurrentGameStatus() types.GameStatus {
	copiedStats := make(map[string]types.PlayerStatus)
	for player, status := range i.playerStatus {
		copiedStats[player.UserID] = *status
	}
	var spec []string

	for _, spectator := range i.spectator {
		spec = append(spec, spectator.UserID)
	}

	gameStats := types.GameStatus{
		IsTimeRunning: i.timer.IsRunning,
		IsVoting:      i.isVoting,
		IsRoundEnd:    i.isRoundEnd,
		RoundLeft:     i.roundLeft,
		PlayerStats:   copiedStats,
		Spectator:     spec,
	}

	return gameStats
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
