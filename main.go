package main

import (
	"github.com/phaseharry/terminal-chess/internal/chess"
)

func main() {
	game := chess.NewGame()
	game.GetAvailableMoves(chess.Position{Row: 0, Col: 0})
}
