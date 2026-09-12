package main

import (
	"fmt"

	"github.com/phaseharry/terminal-chess/internal/chess"
)

func main() {
	game := chess.NewGame()
	_, err := game.GetAvailableMoves(chess.Position{Row: 0, Col: 0})
	if err != nil {
		fmt.Println(err.Error())
	}
}
