package state

// Board holds the state of the game board
type Board struct {
	// PlayerState and OpponentState might be 2D arrays or a different structure
	// depending on how you want to represent the board
	PlayerState   [10][10]string
	OpponentState [10][10]string
}

// NewBoard returns a new Board
func NewBoard() *Board {
	// Initialize the board to some default state
	return &Board{
		PlayerState:   [10][10]string{},
		OpponentState: [10][10]string{},
	}
}

// UpdatePlayerStates updates the player states
func (b *Board) updatePlayerStates(playerState [10][10]string) {
	b.PlayerState = playerState
}
