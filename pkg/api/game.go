package api

import (
	"errors"
	"fmt"
	"strconv"
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

func (g *Game) FireShot(coord string) (FireResult, int, error) {
	result, err := g.client.Fire(FireData{
		Coord: coord},
	)
	if err != nil {
		return FireResult{}, 0, err
	}
	l := g.MarkOpponent(coord, result)
	return result, l, err
}

// StartGame starts the game
func (g *Game) StartGame(nick, desc, targetNick string, coords []string, botGame bool) {
	_, err := g.client.StartGame(nick, desc, targetNick, coords, botGame)
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

func (g *Game) GetDescription() (GameDescription, error) {
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

func (g *Game) MarkOpponent(shot string, result FireResult) int {
	if shot == "" {
		return 0
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
	return g.state.MarkOpponentBoard(x, y, mark)
}

func (g *Game) UpdatePlayerInfo(name string, description string) {
	g.state.UpdatePlayerInfo(name, description)
}

func (g *Game) GetPlayerInfo() (string, string) {
	return g.state.GetPlayerInfo()
}

func (g *Game) UpdatePlayersDesc(d GameDescription) {
	g.state.UpdatePlayersDesc(d.Desc, d.OppDesc)
}

func (g *Game) GetTopPlayerStats() (TopPlayerStats, error) {
	stats, err := g.client.GetTopPlayerStats()
	if err != nil {
		return TopPlayerStats{}, err
	}
	return stats, nil
}

func (g *Game) MarkPlayerShip(coords string) {
	x, y := mapToState(coords)
	g.state.AddShip(x, y)
}

func (g *Game) GetPlayerCoords() []string {
	states := g.state.GetPlayerBoard()

	var coords []string
	for i, row := range states {
		for j, s := range row {
			if s == state.Ship {
				coords = append(coords, mapFromState(i, j))
			}
		}
	}

	// Ensure coordinates range from A1 to J10
	var fixedCoords []string
	for _, coord := range coords {
		if isValidCoord(coord) {
			fixedCoords = append(fixedCoords, coord)
		}
	}

	return fixedCoords
}

func (g *Game) GetPlayerStats(name string) GameStats {
	stats, err := g.client.GetPlayerStats(name)
	if err != nil {
		fmt.Errorf("error getting player stats: %v", err)
		return GameStats{}
	}
	return stats
}

func (g *Game) GetPlayerLobby() []LobbyPlayer {
	lobby, err := g.client.GetLobbyPlayers()
	if err != nil {
		return []LobbyPlayer{}
	}
	return lobby
}

func (g *Game) ClearState() {
	g.state.ClearState()
}

func (g *Game) UpdateLastGameStatus(status string) {
	g.state.UpdateLastGameStatus(status)
}

func (g *Game) LastGameStatus() string {
	return g.state.LastGameStatus()
}

func (g *Game) AbortGame() {
	g.client.AbortGame()
}

func mapFromState(x, y int) string {
	return string(byte(x+65)) + strconv.Itoa(y+1)
}

func isValidCoord(coord string) bool {
	if len(coord) < 2 || len(coord) > 3 {
		return false
	}

	column := coord[0]
	row := coord[1:]
	if column < 'A' || column > 'J' {
		return false
	}

	num, err := strconv.Atoi(row)
	if err != nil || num < 1 || num > 10 {
		return false
	}

	return true
}
