package game

import (
	"context"
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
	defer cancel()

	a.wg.Add(9)

	a.game.StartGame()
	board, err := a.game.LoadPlayerBoard()
	if err != nil {
		return
	}
	_, err = a.game.SetPlayerBoard(board.Board)

	if err != nil {
		return
	}
	go func() {
		defer a.wg.Done()
		a.updateGameStatus(ctx)
	}()

	go func() {
		defer a.wg.Done()
		a.gui.handleGameState(ctx, a.gameStateChannel)
	}()
	go func() {
		defer a.wg.Done()
		a.updateGameState(ctx)
	}()

	go func() {
		defer a.wg.Done()
		a.gui.displayBoard(ctx)
	}()

	go func() {
		defer a.wg.Done()
		a.gui.handleGameStatus(ctx, a.gameStatusChannel, a.playerShotsChannel)
	}()

	go func() {
		defer a.wg.Done()
		a.handleError(ctx, cancel)
	}()

	go func() {
		defer a.wg.Done()
		a.readPlayerShots(ctx)
	}()

	go func() {
		defer a.wg.Done()
		a.gui.listenPlayerShots(ctx, a.playerShotsChannel)
	}()

	a.gui.gui.Start(ctx, nil)

	a.wg.Wait() // Wait for all goroutines to finish
}

// updates game status from the server
func (a *App) updateGameState(ctx context.Context) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			// If the context is canceled, stop the loop
			return
		case <-ticker.C:
			state, err := a.game.GetGameStatus()
			oppShots := state.OppShots
			a.game.MarkOpponentShots(oppShots)
			if err != nil {
				a.errChan <- err // Send error to errChan
				continue
			}
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
	for {
		select {
		case <-ctx.Done():
			// If the context is cancelled, stop the loop
			return
		case err := <-a.errChan:
			a.gui.gui.Log("Error: %v", err)
		}
	}
}

func (a *App) readPlayerShots(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case shot := <-a.playerShotsChannel:
			_, err := a.game.FireShot(shot)
			if err != nil {
				a.errChan <- err
				continue
			}
		}
	}
}
