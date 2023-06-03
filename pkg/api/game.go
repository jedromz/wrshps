package api

import (
	"errors"
	"warships/pkg/state"
)

var (
	ErrAlreadyHit = errors.New("already hit")
)

type Game struct {
	client *Client
	state  *state.GameState
}

// NewGame returns a new Game
func NewGame() *Game {
	return &Game{
		client: NewClient("https://go-pjatk-server.fly.dev/api", ""),
		state:  state.NewGameState(),
	}
}

func (g *Game) FireShot(coord string) (FireResult, error) {
	result, err := g.client.Fire(FireData{
		Coord: coord},
	)
	if err != nil {
		return FireResult{}, err
	}
	g.MarkOpponent(coord, result)
	return result, err
}

// StartGame starts the game
func (g *Game) StartGame() {
	_, err := g.client.StartGame()
	if err != nil {
		return
	}
}
func (g *Game) GetGameStatus() (GameStatus, error) {
	gameState, err := g.client.GetGameStatus()
	if err != nil {
		return GameStatus{}, err
	}

	return gameState, nil
}
func (g *Game) SetPlayerBoard(coords []string) ([10][10]string, error) {
	board, err := g.state.UpdatePlayerBoard(setStatesFromCoords(coords, state.Ship))
	if err != nil {
		return [10][10]string{}, err

	}
	return board, nil
}

func (g *Game) GetDescription() (*GameDescription, error) {
	return g.client.GetGameDescription()
}

func (g *Game) LoadPlayerBoard() (*GameBoard, error) {
	return g.client.GetGameBoard()
}

func (g *Game) UpdateGameState(nick string, desc string, opponent string, oppDesc string) {
	g.state.UpdateGameState(nick, desc, opponent, oppDesc)
}
func (g *Game) GetPlayerBoard() [10][10]string {
	return g.state.GetPlayerBoard()
}
func (g *Game) MarkOpponentShots(shots []string) {
	for _, coord := range shots {
		x, y := mapToState(coord)
		g.state.MarkPlayerBoard(x, y)
	}
}

func (g *Game) GetGameState() (*state.GameState, error) {
	return g.state.GetGameState(), nil
}

func (g *Game) GetOpponentBoard() [10][10]string {
	return g.state.GetOpponentBoard()

}

func (g *Game) MarkOpponent(shot string, result FireResult) {
	if shot == "" {
		return
	}
	x, y := mapToState(shot)
	var mark string
	switch result.Result {
	case "sunk":
		mark = state.Sunk
	case "hit":
		mark = state.Hit
	case "miss":
		mark = state.Miss
	}
	g.state.IncreaseHits(result.Result)
	g.state.MarkOpponentBoard(x, y, mark)
}
