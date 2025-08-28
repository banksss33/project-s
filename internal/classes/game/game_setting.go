package game

type GameSetting struct {
	RoundCount int // round count before the game end | default as 7
	Timer      int // seconds | default as 420 seconds -> 7 min
	Spies      int // can only be 1 or 2 | default as 1
}