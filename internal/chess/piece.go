package chess

type Color string

const (
	White Color = "white"
	Black Color = "black"
)

type PieceType string

const (
	Pawn PieceType = "pawn"
	Knight PieceType = "knight"
	Bishop PieceType = "bishop"
	Rook PieceType = "rook"
	Queen PieceType = "queen"
	King PieceType = "king"
)

type Piece struct {
	pieceType PieceType
	color     Color
}

func NewPiece(color Color, pieceType PieceType) Piece {
	return Piece{
		pieceType ,
		color,
	}
}
