package chess

import (
	"cmp"
	"slices"
	"testing"
)

func TestNewGame(t *testing.T) {
	game := NewGame()
	t.Run("white has the first turn", func(t *testing.T) {
		if game.turn != White {
			t.Errorf("game turn: %s, want %s", game.turn, White)
		}
	})
	t.Run("move history should be an empty slice", func(t *testing.T) {
		if len(game.moves) != 0 {
			t.Errorf("game has move length: %d, want %d", len(game.moves), 0)
		}
	})
	t.Run("castling states should have both kings and rooks has not moved", func(t *testing.T) {
		tests := []struct {
			name  string
			moved bool
		}{
			{"white king", game.castling.whiteKingMoved},
			{"black king", game.castling.blackKingMoved},
			{"white king-side rook", game.castling.whiteKingSideRook.Moved},
			{"white queen-side rook", game.castling.whiteQueenSideRook.Moved},
			{"black king-side rook", game.castling.blackKingSideRook.Moved},
			{"black queen-side rook", game.castling.blackQueenSideRook.Moved},
		}
		for _, tt := range tests {
			if tt.moved {
				t.Errorf("%s Moved = true, want false", tt.name)
			}
		}
	})
}

// func newEmptyGame() *Game {
// 	g := NewGame()
// 	g.board.grid = [8][8]*Piece{}
// 	g.moves = nil
// 	return g
// }

func sortPositions(positions []Position) {
	slices.SortFunc(positions, func(a, b Position) int {
		return cmp.Or(cmp.Compare(a.Row, b.Row), cmp.Compare(a.Col, b.Col))
	})
}

func TestGame_GetAvailableMoves(t *testing.T) {
	type GetAvailableMovesTests struct {
		name    string
		setup   func() *Game
		pos     Position
		want    []Position
		wantErr error
	}

	tests := []GetAvailableMovesTests{
		{
			name:    "Invalid position produces an invalid position error",
			setup:   NewGame,
			pos:     Position{-1, -1},
			wantErr: ErrInvalidPosition,
		},
		{
			name:    "Empty position produces a no piece at position error",
			setup:   NewGame,
			pos:     Position{3, 3},
			wantErr: ErrNoPieceAtPosition,
		},
		// white back rank: every piece is hemmed in by its own pawns except the knights
		{
			name:  "White queen-side rook is blocked at the start",
			setup: NewGame,
			pos:   Position{0, 0},
			want:  []Position{},
		},
		{
			name:  "White queen-side knight can jump the back rank",
			setup: NewGame,
			pos:   Position{0, 1},
			want: []Position{
				{2, 0},
				{2, 2},
			},
		},
		{
			name:  "White queen-side bishop is blocked at the start",
			setup: NewGame,
			pos:   Position{0, 2},
			want:  []Position{},
		},
		{
			name:  "White queen is blocked at the start",
			setup: NewGame,
			pos:   Position{0, 3},
			want:  []Position{},
		},
		{
			name:  "White king is blocked and cannot castle at the start",
			setup: NewGame,
			pos:   Position{0, 4},
			want:  []Position{},
		},
		{
			name:  "White king-side bishop is blocked at the start",
			setup: NewGame,
			pos:   Position{0, 5},
			want:  []Position{},
		},
		{
			name:  "White king-side knight can jump the back rank",
			setup: NewGame,
			pos:   Position{0, 6},
			want: []Position{
				{2, 5},
				{2, 7},
			},
		},
		{
			name:  "White king-side rook is blocked at the start",
			setup: NewGame,
			pos:   Position{0, 7},
			want:  []Position{},
		},
		// white pawns advance one or two squares, with no captures available
		{
			name:  "White pawn (col 0) can advance one or two squares",
			setup: NewGame,
			pos:   Position{1, 0},
			want: []Position{
				{2, 0},
				{3, 0},
			},
		},
		{
			name:  "White pawn (col 1) can advance one or two squares",
			setup: NewGame,
			pos:   Position{1, 1},
			want: []Position{
				{2, 1},
				{3, 1},
			},
		},
		{
			name:  "White pawn (col 2) can advance one or two squares",
			setup: NewGame,
			pos:   Position{1, 2},
			want: []Position{
				{2, 2},
				{3, 2},
			},
		},
		{
			name:  "White pawn (col 3) can advance one or two squares",
			setup: NewGame,
			pos:   Position{1, 3},
			want: []Position{
				{2, 3},
				{3, 3},
			},
		},
		{
			name:  "White pawn (col 4) can advance one or two squares",
			setup: NewGame,
			pos:   Position{1, 4},
			want: []Position{
				{2, 4},
				{3, 4},
			},
		},
		{
			name:  "White pawn (col 5) can advance one or two squares",
			setup: NewGame,
			pos:   Position{1, 5},
			want: []Position{
				{2, 5},
				{3, 5},
			},
		},
		{
			name:  "White pawn (col 6) can advance one or two squares",
			setup: NewGame,
			pos:   Position{1, 6},
			want: []Position{
				{2, 6},
				{3, 6},
			},
		},
		{
			name:  "White pawn (col 7) can advance one or two squares",
			setup: NewGame,
			pos:   Position{1, 7},
			want: []Position{
				{2, 7},
				{3, 7},
			},
		},
		// black pawns advance in the opposite direction
		{
			name:  "Black pawn (col 0) can advance one or two squares",
			setup: NewGame,
			pos:   Position{6, 0},
			want: []Position{
				{5, 0},
				{4, 0},
			},
		},
		{
			name:  "Black pawn (col 1) can advance one or two squares",
			setup: NewGame,
			pos:   Position{6, 1},
			want: []Position{
				{5, 1},
				{4, 1},
			},
		},
		{
			name:  "Black pawn (col 2) can advance one or two squares",
			setup: NewGame,
			pos:   Position{6, 2},
			want: []Position{
				{5, 2},
				{4, 2},
			},
		},
		{
			name:  "Black pawn (col 3) can advance one or two squares",
			setup: NewGame,
			pos:   Position{6, 3},
			want: []Position{
				{5, 3},
				{4, 3},
			},
		},
		{
			name:  "Black pawn (col 4) can advance one or two squares",
			setup: NewGame,
			pos:   Position{6, 4},
			want: []Position{
				{5, 4},
				{4, 4},
			},
		},
		{
			name:  "Black pawn (col 5) can advance one or two squares",
			setup: NewGame,
			pos:   Position{6, 5},
			want: []Position{
				{5, 5},
				{4, 5},
			},
		},
		{
			name:  "Black pawn (col 6) can advance one or two squares",
			setup: NewGame,
			pos:   Position{6, 6},
			want: []Position{
				{5, 6},
				{4, 6},
			},
		},
		{
			name:  "Black pawn (col 7) can advance one or two squares",
			setup: NewGame,
			pos:   Position{6, 7},
			want: []Position{
				{5, 7},
				{4, 7},
			},
		},
		// black back rank mirrors white
		{
			name:  "Black queen-side rook is blocked at the start",
			setup: NewGame,
			pos:   Position{7, 0},
			want:  []Position{},
		},
		{
			name:  "Black queen-side knight can jump the back rank",
			setup: NewGame,
			pos:   Position{7, 1},
			want: []Position{
				{5, 0},
				{5, 2},
			},
		},
		{
			name:  "Black queen-side bishop is blocked at the start",
			setup: NewGame,
			pos:   Position{7, 2},
			want:  []Position{},
		},
		{
			name:  "Black queen is blocked at the start",
			setup: NewGame,
			pos:   Position{7, 3},
			want:  []Position{},
		},
		{
			name:  "Black king is blocked and cannot castle at the start",
			setup: NewGame,
			pos:   Position{7, 4},
			want:  []Position{},
		},
		{
			name:  "Black king-side bishop is blocked at the start",
			setup: NewGame,
			pos:   Position{7, 5},
			want:  []Position{},
		},
		{
			name:  "Black king-side knight can jump the back rank",
			setup: NewGame,
			pos:   Position{7, 6},
			want: []Position{
				{5, 5},
				{5, 7},
			},
		},
		{
			name:  "Black king-side rook is blocked at the start",
			setup: NewGame,
			pos:   Position{7, 7},
			want:  []Position{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := tt.setup()
			got, err := game.GetAvailableMoves(tt.pos)
			if tt.wantErr != err {
				t.Errorf("want error: %v, got error: %v", tt.wantErr, err)
			}
			sortPositions(tt.want)
			sortPositions(got)
			if !slices.Equal(got, tt.want) {
				t.Errorf("GetAvailableMoves(%v) = %v, want %v", tt.pos, got, tt.want)
			}
		})
	}
}
