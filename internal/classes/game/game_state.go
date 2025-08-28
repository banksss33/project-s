package game

type ActionEventType int

const (
	StartTimer ActionEventType = iota
	CountDown
	ShutDown
)

type GameState struct {
	TimeLeft int
}