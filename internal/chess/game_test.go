package chess

import "testing"

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
