package chess

import (
	"fmt"
	"testing"
)

func TestNewBoard(t *testing.T) {
	b := NewBoard()

	tests := []struct {
		name      string
		row, col  int
		color     Color
		pieceType PieceType
	}{
		{
			name:      "white rook one",
			row:       0,
			col:       0,
			color:     White,
			pieceType: Rook,
		},
		{
			name:      "white rook two",
			row:       0,
			col:       7,
			color:     White,
			pieceType: Rook,
		},
		{
			name:      "white knight one",
			row:       0,
			col:       1,
			color:     White,
			pieceType: Knight,
		},
		{
			name:      "white knight two",
			row:       0,
			col:       6,
			color:     White,
			pieceType: Knight,
		},
		{
			name:      "white bishop one",
			row:       0,
			col:       2,
			color:     White,
			pieceType: Bishop,
		},
		{
			name:      "white bishop two",
			row:       0,
			col:       5,
			color:     White,
			pieceType: Bishop,
		},
		{
			name:      "white queen",
			row:       0,
			col:       3,
			color:     White,
			pieceType: Queen,
		},
		{
			name:      "white king",
			row:       0,
			col:       4,
			color:     White,
			pieceType: King,
		},
		{
			name:      "black rook one",
			row:       7,
			col:       0,
			color:     Black,
			pieceType: Rook,
		},
		{
			name:      "black rook two",
			row:       7,
			col:       7,
			color:     Black,
			pieceType: Rook,
		},
		{
			name:      "black knight one",
			row:       7,
			col:       1,
			color:     Black,
			pieceType: Knight,
		},
		{
			name:      "black knight two",
			row:       7,
			col:       6,
			color:     Black,
			pieceType: Knight,
		},
		{
			name:      "black bishop one",
			row:       7,
			col:       2,
			color:     Black,
			pieceType: Bishop,
		},
		{
			name:      "black bishop two",
			row:       7,
			col:       5,
			color:     Black,
			pieceType: Bishop,
		},
		{
			name:      "black queen",
			row:       7,
			col:       3,
			color:     Black,
			pieceType: Queen,
		},
		{
			name:      "black king",
			row:       7,
			col:       4,
			color:     Black,
			pieceType: King,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			piece := b.grid[tt.row][tt.col]
			if piece == nil {
				t.Errorf("grid[%d][%d] is empty, want %s %s", tt.row, tt.col, tt.color, tt.pieceType)
			} else if piece.Color != tt.color || piece.PieceType != tt.pieceType {
				t.Errorf("grid[%d][%d] = %s %s, want %s %s", tt.row, tt.col, piece.Color, piece.PieceType, tt.color, tt.pieceType)
			}
		})
	}

	whitePawnRow := 1
	t.Run("white pawns", func(t *testing.T) {
		for col := range len(b.grid[whitePawnRow]) {
			piece := b.grid[whitePawnRow][col]
			if piece == nil {
				t.Errorf("grid[%d][%d] is empty, want %s %s", whitePawnRow, col, White, Pawn)
			} else if piece.Color != White || piece.PieceType != Pawn {
				t.Errorf("grid[%d][%d] has %s %s, want %s %s", whitePawnRow, col, piece.Color, piece.PieceType, White, Pawn)
			}
		}
	})

	blackPawnRow := 6
	t.Run("black pawns", func(t *testing.T) {
		for col := range len(b.grid[blackPawnRow]) {
			piece := b.grid[blackPawnRow][col]
			if piece == nil {
				t.Errorf("grid[%d][%d] is empty, want %s %s", blackPawnRow, col, Black, Pawn)
			} else if piece.Color != Black || piece.PieceType != Pawn {
				t.Errorf("grid[%d][%d] has %s %s, want %s %s", blackPawnRow, col, piece.Color, piece.PieceType, Black, Pawn)
			}
		}
	})

	t.Run("middle rows are empty", func(t *testing.T) {
		for row := 2; row <= 5; row++ {
			for col := range len(b.grid[row]) {
				piece := b.grid[row][col]
				if piece != nil {
					t.Errorf("grid[%d][%d] is not empty. It has piece %s %s", row, col, piece.Color, piece.PieceType)
				}
			}
		}
	})
}

func TestBoard_IsPositionOutOfBounds(t *testing.T) {
	type IsPositionOutOfBoundTest struct {
		name     string
		position Position
		want     bool
	}
	tests := []IsPositionOutOfBoundTest{
		{"negative row", Position{-1, 0}, true},
		{"negative column", Position{0, -1}, true},
		{"negative row and negative column", Position{-1, -1}, true},
		{"out of bounds row", Position{8, 0}, true},
		{"out of bounds column", Position{0, 8}, true},
	}
	for row := range 8 {
		for col := range 8 {
			tests = append(tests, IsPositionOutOfBoundTest{
				fmt.Sprintf("Row: %d, Column: %d", row, col),
				Position{row, col},
				false,
			})
		}
	}

	b := NewBoard()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := b.IsPositionOutOfBounds(tt.position); got != tt.want {
				t.Errorf("IsPositionOutOfBounds(%v) = %v, want %v", tt.position, got, tt.want)
			}
		})
	}
}
