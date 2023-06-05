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

func (a *App) StartPlayerGame(ctx context.Context) {
	for {
		ctx, cancel := context.WithCancel(ctx)
		var wg sync.WaitGroup

		nick, desc := a.game.GetPlayerInfo()
		fmt.Println("would you like to place your ships?")
		var answer string
		fmt.Scanln(&answer)
		if answer == "y" {
			a.PlaceShips(ctx)
		}
		coords := a.game.GetPlayerCoords()
		fmt.Println(coords)
		var targetNick string
		fmt.Println("Enter target nick: ")
		fmt.Scanln(&targetNick)

		a.game.StartGame(nick, desc, targetNick, coords, false)
		board, err := a.game.LoadPlayerBoard()
		if err != nil {
			a.errChan <- err // Send error to errChan
		}
		_, err = a.game.SetPlayerBoard(board.Board)

		_, err = a.game.GetDescription()

		wg.Add(9)
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
		fmt.Println("Would you like to play again? (y/n)")
		var choice string
		fmt.Scanln(&choice)
		if choice == "n" {
			break
		}
	}

}

func (a *App) StartBotGame(ctx context.Context) {
	for {
		ctx, cancel := context.WithCancel(ctx)

		var wg sync.WaitGroup
		wg.Add(9)
		nick, desc := a.game.GetPlayerInfo()

		fmt.Println("would you like to place your ships?")
		var answer string
		fmt.Scanln(&answer)
		if answer == "y" {
			a.PlaceShips(ctx)
		}
		coords := a.game.GetPlayerCoords()

		a.game.StartGame(nick, desc, "", coords, true)
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
		a.game.LastGameStatus()
		fmt.Println("Would you like to play again? (y/n)")
		var choice string
		fmt.Scanln(&choice)
		if choice == "n" {
			break
		}

	}

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
			a.game.UpdateLastGameStatus(state.LastGameStatus)
			if err != nil {
				a.errChan <- err // Send error to errChan
				continue
			}
			if state.GameStatus == "ended" {
				a.game.ClearState()
				cancel()
				return
			}
			if state.GameStatus == "game_in_progress" {
				d, _ := a.game.GetDescription()
				a.game.UpdatePlayersDesc(d)
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
				PlayerBoard:  state.GetPlayerBoard(),
				OppBoard:     state.GetOpponentBoard(),
				TotalHits:    state.GetTotalHits(),
				TotalShots:   state.GetTotalShots(),
				PlayerDesc:   state.GetPlayerDesc(),
				OppDesc:      state.GetOppDesc(),
				OppShipsSunk: state.GetOppShipsSunk(),
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
			_, _, err := a.game.FireShot(shot)
			if err != nil {
				a.errChan <- err
			}
		}
	}
}

func (a *App) EnterPlayerInfo(ctx context.Context) {
	fmt.Println("Enter your nick: ")
	var name string
	fmt.Scanln(&name)
	fmt.Println("Enter your description: ")
	var description string
	fmt.Scanln(&description)

	a.game.UpdatePlayerInfo(name, description)
}

func (a *App) GetPlayerStats(ctx context.Context) {
	fmt.Println("Enter player nick: ")
	var name string
	fmt.Scanln(&name)
	stats := a.game.GetPlayerStats(name)
	fmt.Println(stats)
}

func (a *App) PrintLobby() {
	lobby := a.game.GetPlayerLobby()
	fmt.Println(lobby)
}
