package game

import (
	"context"
	"fmt"
	gui "github.com/grupawp/warships-gui/v2"
	"strconv"
	"sync"
	"warships/pkg/api"
	"warships/pkg/state"
	// other imports
)

type Gui struct {
	gui            *gui.GUI
	playerBoard    *gui.Board
	opponentBoard  *gui.Board
	playerNick     *gui.Text
	playerDesc     *gui.Text
	opponentNick   *gui.Text
	opponentDesc   *gui.Text
	timer          *gui.Text
	gameStateChan  <-chan *state.GameState
	timerChan      <-chan int
	gameStatusChan chan api.GameStatus
	mu             sync.Mutex
}

func (g *Gui) SetPlayerBoard(states [10][10]gui.State) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.playerBoard.SetStates(states)
}

// NewGui - creates new gui
func NewGui() *Gui {
	return &Gui{
		gui:           gui.NewGUI(true),
		playerNick:    gui.NewText(playerNickX, playerNickY, "Player", nil),
		playerDesc:    gui.NewText(playerDescX, playerDescY, "Your board", nil),
		opponentNick:  gui.NewText(opponentNickX, opponentNickY, "Opponent", nil),
		opponentDesc:  gui.NewText(opponentDescX, opponentDescY, "Opponent board", nil),
		playerBoard:   gui.NewBoard(playerBoardX, playerBoardY, nil),
		opponentBoard: gui.NewBoard(opponentBoardX, opponentBoardY, nil),
		timer:         gui.NewText(timerX, timerY, "", nil),
		mu:            sync.Mutex{},
	}
}
func (g *Gui) setup() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.gui.Draw(g.playerBoard)
	g.gui.Draw(g.opponentBoard)
	g.gui.Draw(g.playerNick)
	g.gui.Draw(g.playerDesc)
	g.gui.Draw(g.opponentNick)
	g.gui.Draw(g.opponentDesc)
	g.gui.Draw(g.timer)
}

func (g *Gui) updateGameState(event GameEvent) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.playerNick.SetText(event.PlayerName)
	g.playerDesc.SetText(event.PlayerDesc)
	g.opponentNick.SetText(event.OpponentName)
	g.opponentDesc.SetText(event.OpponentDesc)
	g.timer.SetText("Timer: " + strconv.Itoa(event.TimeLeft))
}

func (g *Gui) sendPlayerShots(ctx context.Context, shotsChannel chan string) {
	coord := g.playerBoard.Listen(ctx)
	shotsChannel <- coord
}

func (g *Gui) displayBoard(ctx context.Context) {
	g.gui.Draw(g.playerBoard)
	g.gui.Draw(g.opponentBoard)
}

func (g *Gui) handleGameStatus(ctx context.Context, events chan api.GameStatus) {
loop:
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Done handling status OK")
			break loop
		case status := <-events:
			g.mu.Lock()
			g.gui.Draw(gui.NewText(timerX, timerY, "Timer:"+strconv.Itoa(status.Timer), nil))
			g.gui.Draw(gui.NewText(playerNickX, playerNickY, status.Nick, nil))
			g.gui.Draw(gui.NewText(opponentNickX, playerNickY, status.Opponent, nil))
			if status.GameStatus == "ended" {
				g.gui.Draw(gui.NewText(5, 10, "Game ended. Press ctrl + c to go back to the menu", nil))
			}
			g.mu.Unlock()
		}
	}
}

func (g *Gui) handleGameState(ctx context.Context, state chan api.GameState) {
loop:
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Done hanlding state OK")
			break loop
		case gameState := <-state:
			g.mu.Lock()
			g.playerBoard.SetStates(mapStatesToGuiMarks(gameState.PlayerBoard))
			g.opponentBoard.SetStates(mapStatesToGuiMarks(gameState.OppBoard))
			g.gui.Draw(gui.NewText(1, 2, fmt.Sprintf("Accuracy: %s %%",
				getAccuracy(gameState.TotalHits, gameState.TotalShots)), nil))
			g.mu.Unlock()
		}
	}
}

func getAccuracy(hits, shots int) string {
	if shots == 0 {
		return "0.00"
	}
	accuracy := float64(hits) / float64(shots) * 100
	return fmt.Sprintf("%.2f", accuracy)
}

func (g *Gui) listenPlayerShots(ctx context.Context, shots chan string) {
loop:
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Done listening OK")
			break loop
		default:
			shot := g.opponentBoard.Listen(ctx)
			if shot != "" {
				shots <- shot
			}
		}
	}
}

func mapStatesToGuiMarks(sts [10][10]string) [10][10]gui.State {
	var mapped [10][10]gui.State
	for i, row := range sts {
		for j, s := range row {
			if s == state.Sunk {
				s = state.Hit
			}
			mapped[i][j] = gui.State(s)
		}
	}
	return mapped
}
