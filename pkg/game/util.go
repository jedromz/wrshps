package game

import (
	gui "github.com/grupawp/warships-gui/v2"
	"warships/pkg/state"
)

func mapGameStatesToGuiStates(boardStates [10][10]string) [10][10]gui.State {
	var states [10][10]gui.State
	//for each cell in boardStates add appropriate gui.State to states
	for i, row := range boardStates {
		for j, cell := range row {
			switch cell {
			case state.Ship:
				states[i][j] = gui.Ship
			case state.Miss:
				states[i][j] = gui.Miss
			default:
				states[i][j] = gui.Empty
			}
		}
	}
	return states
}
