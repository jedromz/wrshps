// state.go
package state

import (
	"sync"
)

// GameState manages the state of the game
type GameState struct {
	player        *Player
	opponent      *Player
	playerBoard   *Board
	opponentBoard *Board
	// A mutex is used to manage concurrent access to the state
	m sync.Mutex
}

// NewGameState returns a new GameState
func NewGameState() *GameState {
	return &GameState{
		m:             sync.Mutex{},
		player:        &Player{},
		opponent:      &Player{},
		playerBoard:   NewBoard(),
		opponentBoard: NewBoard(),
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

func (g *GameState) MarkOpponentBoard(x int, y int, result string) {
	g.m.Lock()
	defer g.m.Unlock()
	g.opponentBoard.PlayerState[x][y] = result
}

func (g *GameState) IsHitAlready(x, y int) bool {
	g.m.Lock()
	defer g.m.Unlock()
	s := g.opponentBoard.PlayerState[x][y]
	return s == Hit || s == Miss
}
