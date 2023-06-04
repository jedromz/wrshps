package game

import (
	"context"
	"fmt"
	"sync"
	"time"
	"warships/pkg/api"
)

type App struct {
	gui                *Gui
	game               *api.Game
	playerShotsChannel chan string
	gameStatusChannel  chan api.GameStatus
	gameStateChannel   chan api.GameState
	errChan            chan error
	wg                 *sync.WaitGroup // Using a pointer to a WaitGroup
}

func NewApp(gameStatusChannel chan api.GameStatus, playerShotsChannel chan string, gameStateChannel chan api.GameState) *App {
	return &App{
		gui:                NewGui(),
		game:               api.NewGame(),
		playerShotsChannel: playerShotsChannel,
		gameStatusChannel:  gameStatusChannel,
		gameStateChannel:   gameStateChannel,
		errChan:            make(chan error),  // Initializing errChan
		wg:                 &sync.WaitGroup{}, // Initializing the WaitGroup
	}
}

func (a *App) Run(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)

	var wg sync.WaitGroup
	wg.Add(9)

	a.game.StartGame()
	board, err := a.game.LoadPlayerBoard()
	if err != nil {
		a.errChan <- err // Send error to errChan
	}
	_, err = a.game.SetPlayerBoard(board.Board)

	go func() {
		defer wg.Done()
		a.updateGameStatus(ctx)
	}()

	go func() {
		defer wg.Done()
		a.gui.handleGameState(ctx, a.gameStateChannel)
	}()
	go func() {
		defer wg.Done()
		a.updateGameState(ctx, cancel)
	}()

	go func() {
		defer wg.Done()
		a.gui.displayBoard(ctx)
	}()

	go func() {
		defer wg.Done()
		a.gui.handleGameStatus(ctx, a.gameStatusChannel)
	}()

	go func() {
		defer wg.Done()
		a.handleError(ctx, cancel)
	}()

	go func() {
		defer wg.Done()
		a.readPlayerShots(ctx)
	}()

	go func() {
		defer wg.Done()
		a.gui.listenPlayerShots(ctx, a.playerShotsChannel)
	}()
	go func() {
		defer wg.Done()
		a.gui.gui.Start(ctx, nil)
	}()
	wg.Wait() // Wait for all goroutines to finish
	fmt.Println("Game over")

	fmt.Println("What now?")
	fmt.Scanln()

}

// updates game status from the server
func (a *App) updateGameState(ctx context.Context, cancel context.CancelFunc) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
loop:
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Done updating Game state ")
			break loop
		case <-ticker.C:
			state, err := a.game.GetGameStatus()
			if err != nil {
				a.errChan <- err // Send error to errChan
				continue
			}
			if state.GameStatus == "ended" {
				cancel()
				return
			}
			oppShots := state.OppShots
			a.game.MarkOpponentShots(oppShots)
			a.gameStatusChannel <- state
		}
	}
}

// updates game state from the storage
func (a *App) updateGameStatus(ctx context.Context) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Done updating status OK")
			return
		case <-ticker.C:
			state, err := a.game.GetGameState()
			if err != nil {
				a.errChan <- err
				continue
			}
			a.gameStateChannel <- api.GameState{
				PlayerBoard: state.GetPlayerBoard(),
				OppBoard:    state.GetOpponentBoard(),
				TotalHits:   state.GetTotalHits(),
				TotalShots:  state.GetTotalShots(),
			}
		}
	}
}

// Modified handleError to accept a cancel function
func (a *App) handleError(ctx context.Context, cancel context.CancelFunc) {
loop:
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Done handling erros")
			// If the context is cancelled, stop the loop
			break loop
		case err := <-a.errChan:
			a.gui.gui.Log("Error: %v", err)
		}
	}
}

func (a *App) readPlayerShots(ctx context.Context) {
loop:
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Done reading shots OK")
			break loop
		case shot := <-a.playerShotsChannel:
			_, err := a.game.FireShot(shot)
			if err != nil {
				a.errChan <- err
			}
		}
	}
}
