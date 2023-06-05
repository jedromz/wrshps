package game

import (
	"context"
	"fmt"
	gui "github.com/grupawp/warships-gui/v2"
	"sync"
)

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
	invalid := gui.NewText(50, 1, "", nil)
	placeGui := gui.NewGUI(false)
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
					invalid.SetText("")
					currStates = newStates
					fullCoords = append(fullCoords, coords...)
				} else {
					invalid.SetText("Invalid placement, try again")
					placeGui.Draw(invalid)
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
