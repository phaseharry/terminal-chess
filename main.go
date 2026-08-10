package main

import (
	"github.com/phaseharry/terminal-chess/internal/chess"
)

func main() {
	board := chess.NewBoard()
	board.Print()
}
