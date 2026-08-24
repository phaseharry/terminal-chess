package chess

import (
	"fmt"
	"strings"
)

type Board struct {
	grid [8][8]*Piece
}

func NewBoard() *Board {
	b := &Board{}
	b.init()
	return b
}

func (b *Board) init() {
	b.initPieces(Black)
	b.initPieces(White)
}

func (b *Board) initPieces(color Color) {
	// default is white
	pawnRow := 1
	backRank := 0

	if color == Black {
		pawnRow = 6
		backRank = 7
	}

	for i := range 8 {
		pawn := NewPiece(color, Pawn)
		b.grid[pawnRow][i] = &pawn
	}

	knightOne, knightTwo := NewPiece(color, Knight), NewPiece(color, Knight)
	b.grid[backRank][1] = &knightOne
	b.grid[backRank][6] = &knightTwo

	bishopOne, bishopTwo := NewPiece(color, Bishop), NewPiece(color, Bishop)
	b.grid[backRank][2] = &bishopOne
	b.grid[backRank][5] = &bishopTwo

	rookOne, rookTwo := NewPiece(color, Rook), NewPiece(color, Rook)
	b.grid[backRank][0] = &rookOne
	b.grid[backRank][7] = &rookTwo

	queen := NewPiece(color, Queen)
	king := NewPiece(color, King)
	b.grid[backRank][3] = &queen
	b.grid[backRank][4] = &king
}

func (b *Board) Print() {
	for r := range b.grid {
		var row strings.Builder
		row.WriteByte('[')
		for c := range b.grid[r] {
			if c > 0 {
				row.WriteByte(' ')
			}
			piece := b.grid[r][c]
			if piece == nil {
				row.WriteString("  ")
				continue
			}
			if piece.Color == White {
				row.WriteByte('W')
			} else {
				row.WriteByte('B')
			}
			if piece.PieceType == Knight {
				// writing n for knighht so king will still be k
				row.WriteByte('n')
			} else {
				row.WriteByte(piece.PieceType[0])
			}
		}
		row.WriteByte(']')
		fmt.Println(row.String())
	}
}
