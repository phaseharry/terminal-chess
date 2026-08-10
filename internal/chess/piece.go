package chess

import "slices"

type Color string

const (
	White Color = "white"
	Black Color = "black"
)

type PieceType string

const (
	Pawn   PieceType = "pawn"
	Knight PieceType = "knight"
	Bishop PieceType = "bishop"
	Rook   PieceType = "rook"
	Queen  PieceType = "queen"
	King   PieceType = "king"
)

type Piece struct {
	PieceType PieceType
	Color     Color
}

func NewPiece(color Color, pieceType PieceType) Piece {
	return Piece{
		pieceType,
		color,
	}
}

type MoveSpec struct {
	Direction []Moveset
	// if true, then the piece can keep going in the directions listed,
	// otherwise it is only a one jump move
	Sliding bool
}

type Moveset struct {
	DeltaRow int
	DeltaCol int
}

var straightDirections = []Moveset{
	{DeltaRow: 0, DeltaCol: 1},
	{DeltaRow: 1, DeltaCol: 0},
	{DeltaRow: -1, DeltaCol: 0},
	{DeltaRow: 0, DeltaCol: -1},
}

var diagonalDirections = []Moveset{
	{DeltaRow: -1, DeltaCol: 1},
	{DeltaRow: 1, DeltaCol: -1},
	{DeltaRow: 1, DeltaCol: 1},
	{DeltaRow: -1, DeltaCol: -1},
}

var allDirections = slices.Concat(straightDirections, diagonalDirections)

var knightMoveset = []Moveset{
	{DeltaRow: -2, DeltaCol: 1},
	{DeltaRow: -2, DeltaCol: -1},
	{DeltaRow: -1, DeltaCol: 2},
	{DeltaRow: -1, DeltaCol: -2},
	{DeltaRow: 2, DeltaCol: 1},
	{DeltaRow: 2, DeltaCol: -1},
	{DeltaRow: 1, DeltaCol: 2},
	{DeltaRow: 1, DeltaCol: -2},
}

var PieceMoves = map[PieceType]MoveSpec{
	King: {
		Direction: allDirections,
		Sliding:   false,
	},
	Queen: {
		Direction: allDirections,
		Sliding:   true,
	},
	Rook: {
		Direction: straightDirections,
		Sliding:   true,
	},
	Bishop: {
		Direction: diagonalDirections,
		Sliding:   true,
	},
	Knight: {
		Direction: knightMoveset,
		Sliding:   false,
	},
}
