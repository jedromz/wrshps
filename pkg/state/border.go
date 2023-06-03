package state

func (b *Board) DrawBorder(x, y int) [][]int {
	vec := [][]int{
		{1, 0},
		{0, 1},
		{-1, 0},
		{0, -1},
		{1, 1},
		{-1, 1},
		{-1, -1},
		{1, -1},
	}
	shipFound := b.FindShip(x, y)
	for _, v := range vec {
		for _, s := range shipFound {
			xA := s[0] + v[0]
			yA := s[1] + v[1]
			if !isInRange(xA, yA) {
				continue
			}
			if !isShip(xA, yA, b) {
				b.Mark(xA, yA, Miss)
			}
		}
	}
	return shipFound
}

func (b *Board) FindShip(x, y int) [][]int {
	vec := [][]int{
		{1, 0},
		{0, 1},
		{-1, 0},
		{0, -1},
		{1, 1},
		{-1, 1},
		{-1, -1},
		{1, -1},
	}
	shipPlacement := [][]int{
		{x, y},
	}
	for _, v := range vec {
		shipPlacement = append(shipPlacement, findShipRecursive(x, y, v, b)...)
	}
	return shipPlacement
}
func findShipRecursive(x, y int, v []int, b *Board) [][]int {
	if x+v[0] < 0 || x+v[0] >= 10 || y+v[1] < 0 || y+v[1] >= 10 {
		return [][]int{}
	}
	if isShip(x+v[0], y+v[1], b) {
		coords := [][]int{{x + v[0], y + v[1]}}
		recursiveCoords := findShipRecursive(coords[0][0], coords[0][1], v, b)
		return append(coords, recursiveCoords...)
	}
	return [][]int{}
}
func isInRange(x, y int) bool {
	return x >= 0 && x < 10 && y >= 0 && y < 10
}
func isShip(x, y int, b *Board) bool {
	return b.PlayerState[x][y] == Hit || b.PlayerState[x][y] == Sunk
}
