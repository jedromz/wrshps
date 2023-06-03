package game

import (
	"context"
	"errors"
	"fmt"
	"time"
	"warships/pkg/api"
)

type App struct {
	gui        *Gui
	game       *api.Game
	guiChannel chan GameEvent
	timerChan  chan int
	shots      chan string
}

// Setup - sets up ne app
func Setup(c chan GameEvent, s chan string, t chan int) *App {
	return &App{
		gui:        NewGui(),
		game:       api.NewGame(),
		guiChannel: c,
		shots:      s,
		timerChan:  t,
	}
}

func (a *App) Run(ctx context.Context) {
	fmt.Println("Starting game")
	a.game.StartGame()
	a.waitForGameStart(ctx)

	fmt.Println("Game started")
	// get game desc
	desc, err := a.game.GetDescription()
	if err != nil {
		return
	}
	fmt.Println("Updating game description")
	// update playerData in state
	a.game.UpdateGameState(desc.Nick, desc.Desc, desc.Opponent, desc.OppDesc)
	board, err := a.game.LoadPlayerBoard()
	if err != nil {
		return
	}

	fmt.Println("Updating player board")
	a.game.SetPlayerBoard(board.Board)
	playerBoard := a.game.GetPlayerBoard()
	a.gui.SetPlayerBoard(mapGameStatesToGuiStates(playerBoard))
	go a.gui.DisplayBoards(ctx, a.guiChannel, a.shots)
	go a.gameLoop()

	a.gui.StartGui(ctx)
}

func (a *App) gameLoop() bool {

	ticker := time.NewTicker(1000 * time.Millisecond)
	var s api.GameStatus
	var err error
	for s, err = a.game.GetGameStatus(); err == nil && s.GameStatus != "ended"; s, err = a.game.GetGameStatus() {
		a.timerChan <- s.Timer
		a.game.MarkOpponentShots(s.OppShots)
		_, err := a.game.GetGameState()
		if err != nil {
			return true
		}
		a.guiChannel <- GameEvent{
			PlayerStates:   a.game.GetPlayerBoard(),
			OpponentStates: a.game.GetOpponentBoard(),
			PlayerName:     s.Nick,
			PlayerDesc:     "Your board",
			OpponentDesc:   "Opponent board",
			OpponentName:   "Opponent",
			TimeLeft:       s.Timer,
			ShouldFire:     s.ShouldFire,
		}
		if s.ShouldFire {
			shot := <-a.shots
			result, err := a.game.FireShot(shot)
			if err != nil && errors.Is(err, api.ErrAlreadyHit) {
				continue
			}
			a.game.MarkOpponent(shot, result)
		}
		<-ticker.C
	}
	ticker.Stop()
	a.guiChannel <- GameEvent{
		PlayerStates:   a.game.GetPlayerBoard(),
		OpponentStates: a.game.GetOpponentBoard(),
		PlayerName:     "Player",
		PlayerDesc:     "Your board",
		OpponentDesc:   "Opponent board",
		OpponentName:   "Opponent",
		TimeLeft:       0,
		ShouldFire:     false,
		GameState:      "ended",
		Result:         "ed",
	}
	return false
}

func (a *App) waitForGameStart(ctx context.Context) {
	ticker := time.NewTicker(500 * time.Millisecond) // Create a ticker that ticks every second
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			status, err := a.game.GetGameStatus()
			if err != nil {
				return
			}
			if status.GameStatus == "game_in_progress" {
				return
			}
		}
	}
}
