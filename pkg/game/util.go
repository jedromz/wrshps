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
