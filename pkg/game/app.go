package game

import (
	"context"
	"fmt"
	gui "github.com/grupawp/warships-gui/v2"
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

		a.game.StartGame(nick, desc, "", coords, false)
		board, err := a.game.LoadPlayerBoard()
		if err != nil {
			a.errChan <- err // Send error to errChan
		}
		_, err = a.game.SetPlayerBoard(board.Board)

		d, err := a.game.GetDescription()

		a.game.UpdatePlayersDesc(d)
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
				PlayerDesc:  state.GetPlayerDesc(),
				OppDesc:     state.GetOppDesc(),
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

func (a *App) EnterPlayerInfo(ctx context.Context) {
	fmt.Println("Enter your nick: ")
	var name string
	fmt.Scanln(&name)
	fmt.Println("Enter your description: ")
	var description string
	fmt.Scanln(&description)

	a.game.UpdatePlayerInfo(name, description)
}

func (a *App) PlaceShips(ctx context.Context) {
	var mutex sync.Mutex

	ships := map[int]int{
		4: 1,
		3: 2,
		2: 3,
		1: 4,
	}
	currStates := [10][10]gui.State{}
	newStates := [10][10]gui.State{}
	var fullCoords []string

	board := gui.NewBoard(0, 0, nil)
	hint := gui.NewText(50, 0, "Place some ship(s)", nil)
	placeGui := gui.NewGUI(true)
	placeGui.Draw(board)
	placeGui.Draw(hint)
	go func() {
		for k, v := range ships {
			hint.SetText(fmt.Sprintf("Place %v ship(s) of length %v", v, k))
			placeGui.Draw(hint)
			for i := 0; i < v; i++ {
				var coords []string
				for j := 0; j < k; j++ {
					coord := board.Listen(ctx)
					coords = append(coords, coord)
					x, y := mapToState(coord)

					mutex.Lock()
					newStates[x][y] = gui.Ship
					mutex.Unlock()
					board.SetStates(newStates)
				}
				mutex.Lock()
				valid := isValidPlacement(coords) && touchesAnotherShip(coords, currStates) == false
				mutex.Unlock()
				if valid {
					currStates = newStates
					fullCoords = append(fullCoords, coords...)
				} else {
					hint.SetText("Invalid placement, try again")
					placeGui.Draw(hint)
					for _, coord := range coords {
						x, y := mapToState(coord)
						newStates[x][y] = gui.Empty
					}
					board.SetStates(newStates)
					i--
				}
			}
		}
		hint.SetText("Done placing ships. Press ctrl+c save and exit")
		a.game.SetPlayerBoard(fullCoords)
	}()
	placeGui.Start(ctx, nil)
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

func mapToState(coord string) (int, int) {
	if len(coord) > 2 {
		return int(coord[0] - 65), 9
	}
	x := int(coord[0] - 65)
	y := int(coord[1] - 49)
	return x, y
}

// check if the placement is correct
func isAlreadyPlaced(x, y int, states [10][10]gui.State) bool {
	return states[x][y] != gui.Empty
}
func isValidPlacement(coords []string) bool {
	if len(coords) == 0 || len(coords) > 4 {
		return false
	}

	x, y := mapToState(coords[0])
	directions := [][]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} // Right, Left, Down, Up

	for i := 1; i < len(coords); i++ {
		xi, yi := mapToState(coords[i])

		if xi != x && yi != y {
			return false // Invalid coordinate placement (not horizontal or vertical)
		}

		found := false
		for _, dir := range directions {
			dx, dy := dir[0], dir[1]
			if xi == x+dx && yi == y+dy {
				found = true
				break
			}
		}

		if !found {
			return false // Invalid coordinate placement (not adjacent)
		}

		x, y = xi, yi
	}

	return true
}
func touchesAnotherShip(coords []string, states [10][10]gui.State) bool {
	for _, v := range coords {
		x, y := mapToState(v)
		if hasShipAround(x, y, states) {
			return true
		}
	}

	return false
}
func hasShipAround(x, y int, states [10][10]gui.State) bool {
	return isShip(x-1, y, states) || // Left
		isShip(x+1, y, states) || // Right
		isShip(x, y-1, states) || // Up
		isShip(x, y+1, states) || // Down
		isShip(x-1, y-1, states) || // Top Left
		isShip(x-1, y+1, states) || // Bottom Left
		isShip(x+1, y-1, states) || // Top Right
		isShip(x+1, y+1, states) // Bottom Right
}

func isShip(x, y int, states [10][10]gui.State) bool {
	if x > 9 || y > 9 || x < 0 || y < 0 {
		return false
	}
	return states[x][y] == gui.Ship
}
