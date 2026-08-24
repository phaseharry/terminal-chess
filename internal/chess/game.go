package chess

type Game struct {
	board    *Board
	turn     Color
	// castling castlingRights
}

// type castlingRights struct {
// 	white castlingGates
// 	black castlingGates
// }

// type castlingGates struct {
// 	rookOneMoved bool
// 	rookTwoMoved bool
// 	kingMoved    bool
// }

func NewGame() *Game {
	g := &Game{}
	g.initGame()

	return g
}

func (g *Game) initGame() {
	board := NewBoard()
	g.board = board
	g.turn = White
}

func (g *Game) Move() error {
	return nil
}

// func (g *Game) isMoveLegal() bool {
// 	return false
// }

// func (g *Game) isCheck() bool {
// 	return false
// }
