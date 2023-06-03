package game

type GameEvent struct {
	PlayerStates   [10][10]string
	OpponentStates [10][10]string
	PlayerName     string
	PlayerDesc     string
	OpponentName   string
	OpponentDesc   string
	TimeLeft       int
	ShouldFire     bool
	GameState      string
	Result         string
}
