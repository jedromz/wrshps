package game

import (
	"context"
	"fmt"
	gui "github.com/grupawp/warships-gui/v2"
	"strconv"
	"sync"
	"warships/pkg/api"
	"warships/pkg/state"
)

type Gui struct {
	gui            *gui.GUI
	playerBoard    *gui.Board
	opponentBoard  *gui.Board
	playerNick     *gui.Text
	playerDesc     *gui.Text
	opponentNick   *gui.Text
	opponentDesc   *gui.Text
	turn           *gui.Text
	timer          *gui.Text
	waiting        *gui.Text
	numberOf1Ships *gui.Text
	numberOf2Ships *gui.Text
	numberOf3Ships *gui.Text
	numberOf4Ships *gui.Text
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
		gui:            gui.NewGUI(false),
		playerNick:     gui.NewText(playerNickX, playerNickY, "Player", nil),
		playerDesc:     gui.NewText(playerDescX, playerDescY, "Your board", nil),
		opponentNick:   gui.NewText(opponentNickX, opponentNickY, "Opponent", nil),
		opponentDesc:   gui.NewText(opponentDescX, opponentDescY, "Opponent board", nil),
		playerBoard:    gui.NewBoard(playerBoardX, playerBoardY, nil),
		opponentBoard:  gui.NewBoard(opponentBoardX, opponentBoardY, nil),
		waiting:        gui.NewText(10, 10, "Waiting for opponent...", nil),
		turn:           gui.NewText(1, 3, "", nil),
		timer:          gui.NewText(timerX, timerY, "", nil),
		mu:             sync.Mutex{},
		numberOf1Ships: gui.NewText(100, 9, "4 ships of length 1", nil),
		numberOf2Ships: gui.NewText(100, 10, "3 ships of length 2", nil),
		numberOf3Ships: gui.NewText(100, 11, "2 ships of length 3", nil),
		numberOf4Ships: gui.NewText(100, 12, "1 ship of lengthj 4", nil),
	}
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
			break loop
		case status := <-events:
			g.mu.Lock()
			g.updateTurn(status)
			g.updateTimer(status)
			g.updatePlayers(status)
			g.gui.Draw(g.waiting)
			g.gui.Draw(g.numberOf1Ships)
			g.gui.Draw(g.numberOf2Ships)
			g.gui.Draw(g.numberOf3Ships)
			g.gui.Draw(g.numberOf4Ships)
			if status.GameStatus == "ended" {
				g.gui.Draw(gui.NewText(5, 10, "Game ended. Press ctrl + c to go back to the menu", nil))
			}
			if status.GameStatus == "game_in_progress" {
				g.waiting.SetText("")
				g.gui.Draw(g.waiting)
			}
		}
		g.mu.Unlock()
	}
}

func (g *Gui) updatePlayers(status api.GameStatus) {
	g.playerNick.SetText(status.Nick)
	g.gui.Draw(g.playerNick)
	g.opponentNick.SetText(status.Opponent)
	g.gui.Draw(g.opponentNick)
}

func (g *Gui) updateTimer(status api.GameStatus) {
	g.timer.SetText(fmt.Sprintf("Timer: %d", status.Timer))
	g.gui.Draw(g.timer)
}

func (g *Gui) updateTurn(status api.GameStatus) {
	if status.ShouldFire {
		g.turn.SetText("Your turn")
		g.gui.Draw(g.turn)
	} else {
		g.turn.SetText("Opponent's turn")
		g.gui.Draw(g.turn)
	}
}

func (g *Gui) handleGameState(ctx context.Context, state chan api.GameState) {
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case gameState := <-state:
			g.mu.Lock()
			g.playerBoard.SetStates(mapStatesToGuiMarks(gameState.PlayerBoard))
			g.opponentBoard.SetStates(mapStatesToGuiMarks(gameState.OppBoard))
			g.opponentBoard.SetStates(mapStatesToGuiMarks(gameState.OppBoard))
			g.gui.Draw(gui.NewText(playerDescX, playerDescY, gameState.PlayerDesc, nil))
			g.gui.Draw(gui.NewText(opponentDescX, opponentDescY, gameState.OppDesc, nil))
			g.gui.Draw(gui.NewText(1, 2, fmt.Sprintf("Accuracy: %s %%",
				getAccuracy(gameState.TotalHits, gameState.TotalShots)), nil))
			g.mu.Unlock()

			g.numberOf1Ships.SetText(strconv.Itoa(gameState.OppShipsSunk[1]) + " ships of length 1")
			g.numberOf2Ships.SetText(strconv.Itoa(gameState.OppShipsSunk[2]) + " ships of length 2")
			g.numberOf3Ships.SetText(strconv.Itoa(gameState.OppShipsSunk[3]) + " ships of length 3")
			g.numberOf4Ships.SetText(strconv.Itoa(gameState.OppShipsSunk[4]) + " ships of length 4")
			g.drawLegend()
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
	var s []string
loop:

	for {

		select {
		case <-ctx.Done():

			break loop
		default:
			shot := g.opponentBoard.Listen(ctx)
			if shot != "" && !contains(s, shot) {
				s = append(s, shot)
				shots <- shot
			}
		}
	}
}
func (g *Gui) drawLegend() {
	g.gui.Draw(gui.NewText(100, 4, "H - Hit", nil))
	g.gui.Draw(gui.NewText(100, 5, "M - Miss", nil))
	g.gui.Draw(gui.NewText(100, 6, "S - Ship", nil))
	g.gui.Draw(gui.NewText(100, 7, "~ - Empty", nil))
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
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
