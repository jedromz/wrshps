package game

import (
	"context"
	gui "github.com/grupawp/warships-gui/v2"
	"strconv"
	"warships/pkg/state"
)

type Gui struct {
	gui           *gui.GUI
	playerBoard   *gui.Board
	opponentBoard *gui.Board
	playerNick    *gui.Text
	playerDesc    *gui.Text
	opponentNick  *gui.Text
	opponentDesc  *gui.Text
	timer         *gui.Text
	gameStateChan <-chan *state.GameState
	timerChan     <-chan int
}

// NewGui - creates new gui
func NewGui() *Gui {
	return &Gui{
		gui:           gui.NewGUI(true),
		playerNick:    gui.NewText(1, 25, "Player", nil),
		playerDesc:    gui.NewText(1, 26, "Your board", nil),
		opponentNick:  gui.NewText(50, 25, "Opponent", nil),
		opponentDesc:  gui.NewText(50, 26, "Opponent board", nil),
		playerBoard:   gui.NewBoard(1, 4, nil),
		opponentBoard: gui.NewBoard(50, 4, nil),
		timer:         gui.NewText(1, 1, "", nil),
	}
}
func (g *Gui) SetPlayerBoard(states [10][10]gui.State) {
	g.playerBoard.SetStates(states)

}

// DisplayBoards - displays boards
func (g *Gui) DisplayBoards(ctx context.Context, c chan GameEvent, shots chan string) {

	ctx, cancel := context.WithCancel(ctx)
	g.setup()
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return
		case s := <-c:
			g.updateGameState(s)
			if s.GameState == "ended" {
				g.gui.Draw(gui.NewText(1, 2, s.Result, nil))
			}
			if s.ShouldFire {
				g.gui.Draw(gui.NewText(1, 2, "Your turn", nil))
				coords := g.opponentBoard.Listen(ctx)
				shots <- coords
			} else {
				g.gui.Draw(gui.NewText(1, 2, "Enemy turn", nil))
			}
		case timeLeft := <-g.timerChan:
			g.timer.SetText(strconv.Itoa(timeLeft))
		}
	}
}

func (g *Gui) setup() {
	g.gui.Draw(g.playerBoard)
	g.gui.Draw(g.opponentBoard)
	g.gui.Draw(g.playerNick)
	g.gui.Draw(g.playerDesc)
	g.gui.Draw(g.opponentNick)
	g.gui.Draw(g.opponentDesc)
	g.gui.Draw(g.timer)
}

func (g *Gui) StartGui(ctx context.Context) {
	g.gui.Start(ctx, nil)
}

func (g *Gui) updateGameState(event GameEvent) {
	g.playerBoard.SetStates(mapBoardToGuiStates(event.PlayerStates))
	g.opponentBoard.SetStates(mapBoardToGuiStates(event.OpponentStates))
	g.playerNick.SetText(event.PlayerName)
	g.playerDesc.SetText(event.PlayerDesc)
	g.opponentNick.SetText(event.OpponentName)
	g.opponentDesc.SetText(event.OpponentDesc)
	g.timer.SetText(strconv.Itoa(event.TimeLeft))

}

func mapBoardToGuiStates(board [10][10]string) [10][10]gui.State {
	var states [10][10]gui.State
	for i, row := range board {
		for j, cell := range row {
			switch cell {
			case "S":
				states[i][j] = gui.Ship
			case "M":
				states[i][j] = gui.Miss
			case "H":
				states[i][j] = gui.Hit
			default:
				states[i][j] = gui.Empty
			}
		}
	}
	return states
}
