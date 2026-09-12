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

func newEmptyGame() *Game {
	g := NewGame()
	g.board.grid = [8][8]*Piece{}
	g.moves = nil
	return g
}

// getDifferentPositionsBetweenGotAndWant reports which of want are absent from got, and which of got were not wanted.
// Both are nil when the two describe the same set.
func getDifferentPositionsBetweenGotAndWant(got, want []Position) ([]Position, []Position) {
	inGot := make(map[Position]bool, len(got))
	for _, position := range got {
		inGot[position] = true
	}
	inWant := make(map[Position]bool, len(want))
	for _, position := range want {
		inWant[position] = true
	}

	var missing []Position
	var unexpected []Position
	for _, position := range want {
		if !inGot[position] {
			missing = append(missing, position)
		}
	}
	for _, position := range got {
		if !inWant[position] {
			unexpected = append(unexpected, position)
		}
	}
	sortPositions(missing)
	sortPositions(unexpected)
	return missing, unexpected
}

func sortPositions(positions []Position) {
	slices.SortFunc(positions, func(a, b Position) int {
		return cmp.Or(cmp.Compare(a.Row, b.Row), cmp.Compare(a.Col, b.Col))
	})
}

func TestGame_GetAvailableMoves(t *testing.T) {
	type GetAvailableMovesTests struct {
		name  string
		setup func() *Game
		pos   Position
		// used for when we want to test for all the moves available for a piece at current position
		want []Position
		// used for when we want to test that a piece does not explicitly does not include a specific move(s)
		notWant []Position
		// used for when we want to test that the outputted moves includes the want moves
		wantAtLeast []Position
		wantErr     error
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
		{
			name: "White king can castle with queen-side rook if both have not been moved",
			setup: func() *Game {
				game := newEmptyGame()
				king, queenSideRook := NewPiece(White, King), NewPiece(White, Rook)
				game.board.grid[WhiteBackRank][4] = &king
				game.board.grid[WhiteBackRank][0] = &queenSideRook
				return game
			},
			pos:         Position{WhiteBackRank, 4},
			wantAtLeast: []Position{{WhiteBackRank, 2}},
		},
		{
			name: "White king can castle with king-side rook if both have not been moved",
			setup: func() *Game {
				game := newEmptyGame()
				king, kingSideRook := NewPiece(White, King), NewPiece(White, Rook)
				game.board.grid[WhiteBackRank][4] = &king
				game.board.grid[WhiteBackRank][7] = &kingSideRook
				return game
			},
			pos:         Position{WhiteBackRank, 4},
			wantAtLeast: []Position{{WhiteBackRank, 6}},
		},
		{
			name: "White king can castle with both king-side & queen-side rooks if all have not been moved",
			setup: func() *Game {
				game := newEmptyGame()
				king, kingSideRook, queenSideRook := NewPiece(White, King), NewPiece(White, Rook), NewPiece(White, Rook)
				game.board.grid[WhiteBackRank][4] = &king
				game.board.grid[WhiteBackRank][7] = &kingSideRook
				game.board.grid[WhiteBackRank][0] = &queenSideRook
				return game
			},
			pos: Position{WhiteBackRank, 4},
			wantAtLeast: []Position{
				{WhiteBackRank, 2},
				{WhiteBackRank, 6},
			},
		},
		{
			name: "White king cannot castle if the king has been moved but the rooks have not been moved",
			setup: func() *Game {
				game := newEmptyGame()
				game.castling.whiteKingMoved = true
				king, kingSideRook, queenSideRook := NewPiece(White, King), NewPiece(White, Rook), NewPiece(White, Rook)
				game.board.grid[3][6] = &king
				game.board.grid[WhiteBackRank][7] = &kingSideRook
				game.board.grid[WhiteBackRank][0] = &queenSideRook
				return game
			},
			pos: Position{3, 6},
			notWant: []Position{
				{WhiteBackRank, 2},
				{WhiteBackRank, 6},
			},
		},
		{
			name: "White king cannot castle if the rooks have been moved",
			setup: func() *Game {
				game := newEmptyGame()
				game.castling.whiteQueenSideRook.Moved = true
				game.castling.whiteKingSideRook.Moved = true
				king, kingSideRook, queenSideRook := NewPiece(White, King), NewPiece(White, Rook), NewPiece(White, Rook)
				game.board.grid[WhiteBackRank][4] = &king
				game.board.grid[1][7] = &kingSideRook
				game.board.grid[1][0] = &queenSideRook
				return game
			},
			pos: Position{WhiteBackRank, 4},
			notWant: []Position{
				{WhiteBackRank, 2},
				{WhiteBackRank, 6},
			},
		},
		{
			name: "White king cannot castle if there are pieces blocking the cells between the rooks and the king",
			setup: func() *Game {
				game := newEmptyGame()
				whiteKing, whiteKingSideRook, whiteQueenSideRook := NewPiece(White, King), NewPiece(White, Rook), NewPiece(White, Rook)
				whiteBishop, blackKnight := NewPiece(White, Bishop), NewPiece(Black, Knight)

				game.board.grid[WhiteBackRank][0] = &whiteQueenSideRook
				game.board.grid[WhiteBackRank][1] = &whiteBishop
				game.board.grid[WhiteBackRank][4] = &whiteKing
				game.board.grid[WhiteBackRank][5] = &blackKnight
				game.board.grid[WhiteBackRank][7] = &whiteKingSideRook

				return game
			},
			pos: Position{WhiteBackRank, 4},
			notWant: []Position{
				{WhiteBackRank, 2},
				{WhiteBackRank, 6},
			},
		},
		{
			name: "Black king can castle with queen-side rook if both have not been moved",
			setup: func() *Game {
				game := newEmptyGame()
				king, queenSideRook := NewPiece(Black, King), NewPiece(Black, Rook)
				game.board.grid[BlackBackRank][4] = &king
				game.board.grid[BlackBackRank][0] = &queenSideRook
				return game
			},
			pos:         Position{BlackBackRank, 4},
			wantAtLeast: []Position{{BlackBackRank, 2}},
		},
		{
			name: "Black king can castle with king-side rook if both have not been moved",
			setup: func() *Game {
				game := newEmptyGame()
				king, kingSideRook := NewPiece(Black, King), NewPiece(Black, Rook)
				game.board.grid[BlackBackRank][4] = &king
				game.board.grid[BlackBackRank][7] = &kingSideRook
				return game
			},
			pos:         Position{BlackBackRank, 4},
			wantAtLeast: []Position{{BlackBackRank, 6}},
		},
		{
			name: "Black king can castle with both king-side & queen-side rooks if all have not been moved",
			setup: func() *Game {
				game := newEmptyGame()
				king, kingSideRook, queenSideRook := NewPiece(Black, King), NewPiece(Black, Rook), NewPiece(Black, Rook)
				game.board.grid[BlackBackRank][4] = &king
				game.board.grid[BlackBackRank][7] = &kingSideRook
				game.board.grid[BlackBackRank][0] = &queenSideRook
				return game
			},
			pos: Position{BlackBackRank, 4},
			wantAtLeast: []Position{
				{BlackBackRank, 2},
				{BlackBackRank, 6},
			},
		},
		{
			name: "Black king cannot castle if the king has been moved but the rooks have not been moved",
			setup: func() *Game {
				game := newEmptyGame()
				game.castling.blackKingMoved = true
				king, kingSideRook, queenSideRook := NewPiece(Black, King), NewPiece(Black, Rook), NewPiece(Black, Rook)
				game.board.grid[4][6] = &king
				game.board.grid[BlackBackRank][7] = &kingSideRook
				game.board.grid[BlackBackRank][0] = &queenSideRook
				return game
			},
			pos: Position{4, 6},
			notWant: []Position{
				{BlackBackRank, 2},
				{BlackBackRank, 6},
			},
		},
		{
			name: "Black king cannot castle if the rooks have been moved",
			setup: func() *Game {
				game := newEmptyGame()
				game.castling.blackQueenSideRook.Moved = true
				game.castling.blackKingSideRook.Moved = true
				king, kingSideRook, queenSideRook := NewPiece(Black, King), NewPiece(Black, Rook), NewPiece(Black, Rook)
				game.board.grid[BlackBackRank][4] = &king
				game.board.grid[6][7] = &kingSideRook
				game.board.grid[6][0] = &queenSideRook
				return game
			},
			pos: Position{BlackBackRank, 4},
			notWant: []Position{
				{BlackBackRank, 2},
				{BlackBackRank, 6},
			},
		},
		{
			name: "Black king cannot castle if there are pieces blocking the cells between the rooks and the king",
			setup: func() *Game {
				game := newEmptyGame()
				blackKing, blackKingSideRook, blackQueenSideRook := NewPiece(Black, King), NewPiece(Black, Rook), NewPiece(Black, Rook)
				blackBishop, whiteKnight := NewPiece(Black, Bishop), NewPiece(White, Knight)

				game.board.grid[BlackBackRank][0] = &blackQueenSideRook
				game.board.grid[BlackBackRank][1] = &blackBishop
				game.board.grid[BlackBackRank][4] = &blackKing
				game.board.grid[BlackBackRank][5] = &whiteKnight
				game.board.grid[BlackBackRank][7] = &blackKingSideRook

				return game
			},
			pos: Position{BlackBackRank, 4},
			notWant: []Position{
				{BlackBackRank, 2},
				{BlackBackRank, 6},
			},
		},
		{
			name: "White pawn can move forward diagonally one space if there is an opposing black piece at that space",
			setup: func() *Game {
				game := newEmptyGame()
				whitePawn, blackQueen := NewPiece(White, Pawn), NewPiece(Black, Queen)
				game.board.grid[4][4] = &whitePawn
				game.board.grid[5][5] = &blackQueen
				return game
			},
			pos: Position{4, 4},
			wantAtLeast: []Position{
				{5, 5},
			},
		},
		{
			name: "White pawn can only move forward one space at a time after it has left its initial starting position",
			setup: func() *Game {
				game := newEmptyGame()
				game.moves = append(game.moves, MoveDetails{
					Piece:          NewPiece(Black, Pawn),
					SourcePosition: Position{BlackBackRank, 4},
					TargetPosition: Position{BlackBackRank - 2, 4},
				})
				whitePawn := NewPiece(White, Pawn)
				game.board.grid[4][4] = &whitePawn
				return game
			},
			pos: Position{4, 4},
			want: []Position{
				{5, 4},
			},
		},
		{
			name: "White pawn can perform an enpassant if the last move was a black pawn that moved 2 spaces directly adjacent to it",
			setup: func() *Game {
				game := newEmptyGame()
				whitePawn, blackPawn := NewPiece(White, Pawn), NewPiece(Black, Pawn)
				game.moves = append(game.moves, MoveDetails{
					Piece:          blackPawn,
					SourcePosition: Position{BlackBackRank, 4},
					TargetPosition: Position{BlackBackRank - 2, 4},
				})
				game.board.grid[5][4] = &blackPawn
				game.board.grid[5][5] = &whitePawn
				return game
			},
			pos: Position{5, 5},
			wantAtLeast: []Position{
				{6, 4},
			},
		},
		{
			name: "Black pawn can move forward diagonally one space if there is an opposing white piece at that space",
			setup: func() *Game {
				game := newEmptyGame()
				blackPawn, whiteQueen := NewPiece(Black, Pawn), NewPiece(White, Queen)
				game.board.grid[4][4] = &blackPawn
				game.board.grid[3][5] = &whiteQueen
				return game
			},
			pos: Position{4, 4},
			wantAtLeast: []Position{
				{3, 5},
			},
		},
		{
			name: "Black pawn can only move forward one space at a time after it has left its initial starting position",
			setup: func() *Game {
				game := newEmptyGame()
				game.moves = append(game.moves, MoveDetails{
					Piece:          NewPiece(White, Pawn),
					SourcePosition: Position{WhiteBackRank, 4},
					TargetPosition: Position{WhiteBackRank + 2, 4},
				})
				blackPawn := NewPiece(Black, Pawn)
				game.board.grid[4][4] = &blackPawn
				return game
			},
			pos: Position{4, 4},
			want: []Position{
				{3, 4},
			},
		},
		{
			name: "Black pawn can perform an enpassant if the last move was a white pawn that moved 2 spaces directly adjacent to it",
			setup: func() *Game {
				game := newEmptyGame()
				blackPawn, whitePawn := NewPiece(Black, Pawn), NewPiece(White, Pawn)
				game.moves = append(game.moves, MoveDetails{
					Piece:          whitePawn,
					SourcePosition: Position{WhiteBackRank, 4},
					TargetPosition: Position{WhiteBackRank + 2, 4},
				})
				game.board.grid[2][4] = &whitePawn
				game.board.grid[2][5] = &blackPawn
				return game
			},
			pos: Position{2, 5},
			wantAtLeast: []Position{
				{1, 4},
			},
		},
		{
			name: "White rook on an open board slides along its whole rank and file",
			setup: func() *Game {
				game := newEmptyGame()
				whiteRook := NewPiece(White, Rook)
				game.board.grid[4][4] = &whiteRook
				return game
			},
			pos: Position{4, 4},
			want: []Position{
				{0, 4},
				{1, 4},
				{2, 4},
				{3, 4},
				{4, 0},
				{4, 1},
				{4, 2},
				{4, 3},
				{4, 5},
				{4, 6},
				{4, 7},
				{5, 4},
				{6, 4},
				{7, 4},
			},
		},
		{
			name: "White rook stops before a friendly piece and captures an enemy one",
			setup: func() *Game {
				game := newEmptyGame()
				whiteRook := NewPiece(White, Rook)
				game.board.grid[4][4] = &whiteRook
				whitePawn := NewPiece(White, Pawn)
				game.board.grid[4][6] = &whitePawn
				blackPawn := NewPiece(Black, Pawn)
				game.board.grid[2][4] = &blackPawn
				return game
			},
			pos: Position{4, 4},
			want: []Position{
				{2, 4},
				{3, 4},
				{4, 0},
				{4, 1},
				{4, 2},
				{4, 3},
				{4, 5},
				{5, 4},
				{6, 4},
				{7, 4},
			},
		},
		{
			name: "White bishop on an open board slides along all four diagonals",
			setup: func() *Game {
				game := newEmptyGame()
				whiteBishop := NewPiece(White, Bishop)
				game.board.grid[4][4] = &whiteBishop
				return game
			},
			pos: Position{4, 4},
			want: []Position{
				{0, 0},
				{1, 1},
				{1, 7},
				{2, 2},
				{2, 6},
				{3, 3},
				{3, 5},
				{5, 3},
				{5, 5},
				{6, 2},
				{6, 6},
				{7, 1},
				{7, 7},
			},
		},
		{
			name: "White bishop stops before a friendly piece and captures an enemy one",
			setup: func() *Game {
				game := newEmptyGame()
				whiteBishop := NewPiece(White, Bishop)
				game.board.grid[4][4] = &whiteBishop
				whitePawn := NewPiece(White, Pawn)
				game.board.grid[6][6] = &whitePawn
				blackPawn := NewPiece(Black, Pawn)
				game.board.grid[2][2] = &blackPawn
				return game
			},
			pos: Position{4, 4},
			want: []Position{
				{1, 7},
				{2, 2},
				{2, 6},
				{3, 3},
				{3, 5},
				{5, 3},
				{5, 5},
				{6, 2},
				{7, 1},
			},
		},
		{
			name: "White knight on an open board reaches all eight squares",
			setup: func() *Game {
				game := newEmptyGame()
				whiteKnight := NewPiece(White, Knight)
				game.board.grid[4][4] = &whiteKnight
				return game
			},
			pos: Position{4, 4},
			want: []Position{
				{2, 3},
				{2, 5},
				{3, 2},
				{3, 6},
				{5, 2},
				{5, 6},
				{6, 3},
				{6, 5},
			},
		},
		{
			name: "White knight jumps over surrounding pieces but not onto a friendly one",
			setup: func() *Game {
				game := newEmptyGame()
				whiteKnight := NewPiece(White, Knight)
				game.board.grid[4][4] = &whiteKnight
				whitePawnOne := NewPiece(White, Pawn)
				game.board.grid[3][3] = &whitePawnOne
				whitePawnTwo := NewPiece(White, Pawn)
				game.board.grid[3][4] = &whitePawnTwo
				whitePawnThree := NewPiece(White, Pawn)
				game.board.grid[3][5] = &whitePawnThree
				whitePawnFour := NewPiece(White, Pawn)
				game.board.grid[4][3] = &whitePawnFour
				whitePawnFive := NewPiece(White, Pawn)
				game.board.grid[4][5] = &whitePawnFive
				whitePawnSix := NewPiece(White, Pawn)
				game.board.grid[5][3] = &whitePawnSix
				whitePawnSeven := NewPiece(White, Pawn)
				game.board.grid[5][4] = &whitePawnSeven
				whitePawnEight := NewPiece(White, Pawn)
				game.board.grid[5][5] = &whitePawnEight
				whiteRook := NewPiece(White, Rook)
				game.board.grid[6][5] = &whiteRook
				blackRook := NewPiece(Black, Rook)
				game.board.grid[6][3] = &blackRook
				return game
			},
			pos: Position{4, 4},
			want: []Position{
				{2, 3},
				{2, 5},
				{3, 2},
				{3, 6},
				{5, 2},
				{5, 6},
				{6, 3},
			},
		},
		{
			name: "White queen on an open board combines rook and bishop movement",
			setup: func() *Game {
				game := newEmptyGame()
				whiteQueen := NewPiece(White, Queen)
				game.board.grid[4][4] = &whiteQueen
				return game
			},
			pos: Position{4, 4},
			want: []Position{
				{0, 0},
				{0, 4},
				{1, 1},
				{1, 4},
				{1, 7},
				{2, 2},
				{2, 4},
				{2, 6},
				{3, 3},
				{3, 4},
				{3, 5},
				{4, 0},
				{4, 1},
				{4, 2},
				{4, 3},
				{4, 5},
				{4, 6},
				{4, 7},
				{5, 3},
				{5, 4},
				{5, 5},
				{6, 2},
				{6, 4},
				{6, 6},
				{7, 1},
				{7, 4},
				{7, 7},
			},
		},
		{
			name: "White king on an open board reaches its eight neighbours",
			setup: func() *Game {
				game := newEmptyGame()
				game.castling.whiteKingMoved = true
				whiteKing := NewPiece(White, King)
				game.board.grid[4][4] = &whiteKing
				return game
			},
			pos: Position{4, 4},
			want: []Position{
				{3, 3},
				{3, 4},
				{3, 5},
				{4, 3},
				{4, 5},
				{5, 3},
				{5, 4},
				{5, 5},
			},
		},
		{
			name: "Black rook on an open board slides along its whole rank and file",
			setup: func() *Game {
				game := newEmptyGame()
				blackRook := NewPiece(Black, Rook)
				game.board.grid[4][4] = &blackRook
				return game
			},
			pos: Position{4, 4},
			want: []Position{
				{0, 4},
				{1, 4},
				{2, 4},
				{3, 4},
				{4, 0},
				{4, 1},
				{4, 2},
				{4, 3},
				{4, 5},
				{4, 6},
				{4, 7},
				{5, 4},
				{6, 4},
				{7, 4},
			},
		},
		{
			name: "Black rook stops before a friendly piece and captures an enemy one",
			setup: func() *Game {
				game := newEmptyGame()
				blackRook := NewPiece(Black, Rook)
				game.board.grid[4][4] = &blackRook
				blackPawn := NewPiece(Black, Pawn)
				game.board.grid[4][6] = &blackPawn
				whitePawn := NewPiece(White, Pawn)
				game.board.grid[2][4] = &whitePawn
				return game
			},
			pos: Position{4, 4},
			want: []Position{
				{2, 4},
				{3, 4},
				{4, 0},
				{4, 1},
				{4, 2},
				{4, 3},
				{4, 5},
				{5, 4},
				{6, 4},
				{7, 4},
			},
		},
		{
			name: "Black bishop on an open board slides along all four diagonals",
			setup: func() *Game {
				game := newEmptyGame()
				blackBishop := NewPiece(Black, Bishop)
				game.board.grid[4][4] = &blackBishop
				return game
			},
			pos: Position{4, 4},
			want: []Position{
				{0, 0},
				{1, 1},
				{1, 7},
				{2, 2},
				{2, 6},
				{3, 3},
				{3, 5},
				{5, 3},
				{5, 5},
				{6, 2},
				{6, 6},
				{7, 1},
				{7, 7},
			},
		},
		{
			name: "Black bishop stops before a friendly piece and captures an enemy one",
			setup: func() *Game {
				game := newEmptyGame()
				blackBishop := NewPiece(Black, Bishop)
				game.board.grid[4][4] = &blackBishop
				blackPawn := NewPiece(Black, Pawn)
				game.board.grid[6][6] = &blackPawn
				whitePawn := NewPiece(White, Pawn)
				game.board.grid[2][2] = &whitePawn
				return game
			},
			pos: Position{4, 4},
			want: []Position{
				{1, 7},
				{2, 2},
				{2, 6},
				{3, 3},
				{3, 5},
				{5, 3},
				{5, 5},
				{6, 2},
				{7, 1},
			},
		},
		{
			name: "Black knight on an open board reaches all eight squares",
			setup: func() *Game {
				game := newEmptyGame()
				blackKnight := NewPiece(Black, Knight)
				game.board.grid[4][4] = &blackKnight
				return game
			},
			pos: Position{4, 4},
			want: []Position{
				{2, 3},
				{2, 5},
				{3, 2},
				{3, 6},
				{5, 2},
				{5, 6},
				{6, 3},
				{6, 5},
			},
		},
		{
			name: "Black knight jumps over surrounding pieces but not onto a friendly one",
			setup: func() *Game {
				game := newEmptyGame()
				blackKnight := NewPiece(Black, Knight)
				game.board.grid[4][4] = &blackKnight
				blackPawnOne := NewPiece(Black, Pawn)
				game.board.grid[3][3] = &blackPawnOne
				blackPawnTwo := NewPiece(Black, Pawn)
				game.board.grid[3][4] = &blackPawnTwo
				blackPawnThree := NewPiece(Black, Pawn)
				game.board.grid[3][5] = &blackPawnThree
				blackPawnFour := NewPiece(Black, Pawn)
				game.board.grid[4][3] = &blackPawnFour
				blackPawnFive := NewPiece(Black, Pawn)
				game.board.grid[4][5] = &blackPawnFive
				blackPawnSix := NewPiece(Black, Pawn)
				game.board.grid[5][3] = &blackPawnSix
				blackPawnSeven := NewPiece(Black, Pawn)
				game.board.grid[5][4] = &blackPawnSeven
				blackPawnEight := NewPiece(Black, Pawn)
				game.board.grid[5][5] = &blackPawnEight
				blackRook := NewPiece(Black, Rook)
				game.board.grid[6][5] = &blackRook
				whiteRook := NewPiece(White, Rook)
				game.board.grid[6][3] = &whiteRook
				return game
			},
			pos: Position{4, 4},
			want: []Position{
				{2, 3},
				{2, 5},
				{3, 2},
				{3, 6},
				{5, 2},
				{5, 6},
				{6, 3},
			},
		},
		{
			name: "Black queen on an open board combines rook and bishop movement",
			setup: func() *Game {
				game := newEmptyGame()
				blackQueen := NewPiece(Black, Queen)
				game.board.grid[4][4] = &blackQueen
				return game
			},
			pos: Position{4, 4},
			want: []Position{
				{0, 0},
				{0, 4},
				{1, 1},
				{1, 4},
				{1, 7},
				{2, 2},
				{2, 4},
				{2, 6},
				{3, 3},
				{3, 4},
				{3, 5},
				{4, 0},
				{4, 1},
				{4, 2},
				{4, 3},
				{4, 5},
				{4, 6},
				{4, 7},
				{5, 3},
				{5, 4},
				{5, 5},
				{6, 2},
				{6, 4},
				{6, 6},
				{7, 1},
				{7, 4},
				{7, 7},
			},
		},
		{
			name: "Black king on an open board reaches its eight neighbours",
			setup: func() *Game {
				game := newEmptyGame()
				game.castling.blackKingMoved = true
				blackKing := NewPiece(Black, King)
				game.board.grid[4][4] = &blackKing
				return game
			},
			pos: Position{4, 4},
			want: []Position{
				{3, 3},
				{3, 4},
				{3, 5},
				{4, 3},
				{4, 5},
				{5, 3},
				{5, 4},
				{5, 5},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := tt.setup()
			got, err := game.GetAvailableMoves(tt.pos)
			if tt.wantErr != err {
				t.Errorf("want error: %v, got error: %v", tt.wantErr, err)
			}

			if tt.wantAtLeast != nil {
				if missing, _ := getDifferentPositionsBetweenGotAndWant(got, tt.wantAtLeast); len(missing) > 0 {
					t.Errorf("GetAvailableMoves(%v) = %v, missing %v", tt.pos, got, missing)
				}
			}

			if tt.notWant != nil {
				for _, position := range tt.notWant {
					if slices.Contains(got, position) {
						t.Errorf("GetAvailableMoves(%v) = %v, should not include %v", tt.pos, got, position)
					}
				}
			}

			if tt.want != nil {
				missing, unexpected := getDifferentPositionsBetweenGotAndWant(got, tt.want)
				if len(missing) > 0 || len(unexpected) > 0 {
					t.Errorf("GetAvailableMoves(%v) = %v: missing %v, unexpected %v", tt.pos, got, missing, unexpected)
				}
			}
		})
	}
}
