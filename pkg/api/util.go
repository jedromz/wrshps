package api

// setStatesFromCoords converts []string to [][]string
func setStatesFromCoords(coords []string, s string) [10][10]string {
	state := [10][10]string{}
	for i := range state {
		state[i] = [10]string{}
	}
	for _, coord := range coords {
		x, y := mapToState(coord)
		state[x][y] = s
	}
	return state
}

// mapToState converts string to int
func mapToState(coord string) (int, int) {
	if len(coord) > 2 {
		return int(coord[0] - 65), 9
	}
	x := int(coord[0] - 65)
	y := int(coord[1] - 49)
	return x, y
}
