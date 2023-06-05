package state

import (
	"sync"
)

// GameState manages the state of the game
type GameState struct {
	player         *Player
	opponent       *Player
	playerBoard    *Board
	opponentBoard  *Board
	totalShots     int
	hits           int
	m              sync.Mutex
	lastGameStatus string
	oppShipsSun    map[int]int
}

// NewGameState returns a new GameState
func NewGameState() *GameState {
	return &GameState{
		m:             sync.Mutex{},
		player:        &Player{},
		opponent:      &Player{},
		playerBoard:   NewBoard(),
		opponentBoard: NewBoard(),
		oppShipsSun: map[int]int{
			1: 4,
			2: 3,
			3: 2,
			4: 1,
		},
	}
}

// GetGameState returns the game state
func (g *GameState) GetGameState() *GameState {
	g.m.Lock()
	defer g.m.Unlock()
	return g
}

// UpdateGameState updates the game state
func (g *GameState) UpdateGameState(nick, desc, opp, oppdesc string) {
	g.m.Lock()
	defer g.m.Unlock()
	g.player.Nick = nick
	g.player.Description = desc
	g.opponent.Nick = opp
	g.opponent.Description = oppdesc
}

// UpdatePlayerBoard updates the player board
func (g *GameState) UpdatePlayerBoard(playerState [10][10]string) ([10][10]string, error) {
	g.m.Lock()
	defer g.m.Unlock()
	g.playerBoard.updatePlayerStates(playerState)
	return g.playerBoard.PlayerState, nil
}
func (g *GameState) UpdateOpponentBoard(opponentState [10][10]string) ([10][10]string, error) {
	g.m.Lock()
	defer g.m.Unlock()
	g.opponentBoard.updatePlayerStates(opponentState)
	return g.opponentBoard.PlayerState, nil
}
func (g *GameState) GetPlayerBoard() [10][10]string {
	g.m.Lock()
	defer g.m.Unlock()
	return g.playerBoard.PlayerState
}
func (g *GameState) MarkPlayerBoard(x, y int) {
	g.m.Lock()
	defer g.m.Unlock()
	switch g.playerBoard.PlayerState[x][y] {
	case Ship:
		g.playerBoard.PlayerState[x][y] = Hit
	case Empty:
		g.playerBoard.PlayerState[x][y] = Miss
	}
}

func (g *GameState) GetOpponentBoard() [10][10]string {
	g.m.Lock()
	defer g.m.Unlock()
	return g.opponentBoard.PlayerState
}

func (g *GameState) MarkOpponentBoard(x int, y int, result string) int {
	g.m.Lock()
	defer g.m.Unlock()
	if result == Sunk {
		g.opponentBoard.PlayerState[x][y] = result
		_, l := g.opponentBoard.DrawBorder(x, y)
		g.oppShipsSun[l]--
		return l
	}
	g.opponentBoard.PlayerState[x][y] = result

	return 0
}

func (g *GameState) IsHitAlready(x, y int) bool {
	g.m.Lock()
	defer g.m.Unlock()
	s := g.opponentBoard.PlayerState[x][y]
	return s == Hit || s == Miss
}

func (g *GameState) IncreaseHits(str string) {
	g.m.Lock()
	defer g.m.Unlock()
	if str == "hit" || str == "sunk" {
		g.hits++
	}
	g.totalShots++
}

func (g *GameState) GetTotalShots() int {
	g.m.Lock()
	defer g.m.Unlock()
	return g.totalShots
}
func (g *GameState) GetTotalHits() int {
	g.m.Lock()
	defer g.m.Unlock()
	return g.hits
}

func (g *GameState) UpdatePlayerInfo(name string, description string) {
	g.m.Lock()
	defer g.m.Unlock()
	g.player.Nick = name
	g.player.Description = description
}

func (g *GameState) GetPlayerInfo() (string, string) {
	g.m.Lock()
	defer g.m.Unlock()
	return g.player.Nick, g.player.Description
}

func (g *GameState) GetOppDesc() string {
	g.m.Lock()
	defer g.m.Unlock()
	return g.opponent.Description
}
func (g *GameState) GetPlayerDesc() string {
	g.m.Lock()
	defer g.m.Unlock()
	return g.player.Description
}

func (g *GameState) UpdatePlayersDesc(desc, oppDesc string) {
	g.m.Lock()
	defer g.m.Unlock()
	g.player.Description = desc
	g.opponent.Description = oppDesc
}

func (g *GameState) AddShip(x int, y int) {
	g.m.Lock()
	defer g.m.Unlock()
	g.playerBoard.PlayerState[x][y] = Ship
}

func (g *GameState) ClearState() {
	g.m.Lock()
	defer g.m.Unlock()
	g.playerBoard = NewBoard()
	g.opponentBoard = NewBoard()
	g.totalShots = 0
	g.hits = 0
}

func (g *GameState) GetOppShipsSunk() map[int]int {
	g.m.Lock()
	defer g.m.Unlock()
	return g.oppShipsSun
}

func (g *GameState) UpdateLastGameStatus(status string) {
	g.m.Lock()
	defer g.m.Unlock()
	g.lastGameStatus = status
}

func (g *GameState) LastGameStatus() string {
	g.m.Lock()
	defer g.m.Unlock()
	return g.lastGameStatus
}
