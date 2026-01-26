package types

import (
	"time"
)

type PlayerStatus struct {
	Score       int
	Roles       string
	AlreadyVote bool
}

type GameSetting struct {
	Round     int
	Spies     int
	Timer     int
	Locations map[string][]string //key: location name| value: location roles
}

type GameTimer struct {
	Tick      *time.Ticker
	IsRunning bool
	Countdown int
}

type GameStatus struct {
	IsTimeRunning bool
	IsVoting      bool
	IsRoundEnd    bool
	RoundLeft     int
	PlayerStats   map[string]PlayerStatus
	Spectator     []string
}
