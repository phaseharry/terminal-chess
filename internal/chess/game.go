package chess

import (
	"errors"
	"math"
)

var ErrInvalidPosition = errors.New("invalid position")

var ErrNoPieceAtPosition = errors.New("no piece at position to move")

type Game struct {
	board    *Board
	turn     Color
	moves    []MoveDetails
	castling castlingState
}

type MoveDetails struct {
	Piece          Piece
	SourcePosition Position
	TargetPosition Position
}

type castlingState struct {
	whiteKingSideRook  rookStatus
	whiteQueenSideRook rookStatus
	whiteKingMoved     bool
	blackKingSideRook  rookStatus
	blackQueenSideRook rookStatus
	blackKingMoved     bool
}

type rookStatus struct {
	Rook  *Piece
	Moved bool
}

func NewGame() *Game {
	g := &Game{}
	g.initGame()

	return g
}

func (g *Game) initGame() {
	board := NewBoard()
	g.board = board
	g.turn = White
	g.castling = castlingState{
		whiteKingSideRook: rookStatus{
			Rook:  g.board.grid[WhiteBackRow][7],
			Moved: false,
		},
		whiteQueenSideRook: rookStatus{
			Rook:  g.board.grid[WhiteBackRow][0],
			Moved: false,
		},
		blackKingMoved: false,
		blackKingSideRook: rookStatus{
			Rook:  g.board.grid[BlackBackRow][7],
			Moved: false,
		},
		blackQueenSideRook: rookStatus{
			Rook:  g.board.grid[BlackBackRow][0],
			Moved: false,
		},
	}
}

func (g *Game) GetAvailableMoves(sourcePosition Position) ([]Position, error) {
	var availablePositions []Position
	sourceRow, sourceCol := sourcePosition.Row, sourcePosition.Col

	// should never happen
	if g.board.IsPositionOutOfBounds(sourcePosition) {
		return availablePositions, ErrInvalidPosition
	}

	if g.board.grid[sourceRow][sourceCol] == nil {
		return availablePositions, ErrNoPieceAtPosition
	}

	piece := g.board.grid[sourceRow][sourceCol]

	switch piece.PieceType {
	case King:
		availablePositions = g.getMovesForPiece(sourcePosition, PieceMoves[piece.PieceType])
		availablePositions = append(availablePositions, g.getMovesForCastlablePieces(sourcePosition)...)
	case Pawn:
		availablePositions = g.getMovesForPawn(sourcePosition)
	default:
		availablePositions = g.getMovesForPiece(sourcePosition, PieceMoves[piece.PieceType])
	}

	return availablePositions, nil
}

func (g *Game) getMovesForCastlablePieces(sourcePosition Position) []Position {
	var castleMoves []Position
	sourcePiece := g.board.grid[sourcePosition.Row][sourcePosition.Col]
	sourceColor := sourcePiece.Color

	backrow := WhiteBackRow
	kingMoved := g.castling.whiteKingMoved
	kingSideRookMoved := g.castling.whiteKingSideRook.Moved
	queenSideRookMoved := g.castling.whiteQueenSideRook.Moved

	if sourceColor == Black {
		backrow = BlackBackRow
		kingMoved = g.castling.blackKingMoved
		kingSideRookMoved = g.castling.blackKingSideRook.Moved
		queenSideRookMoved = g.castling.blackQueenSideRook.Moved
	}

	// king has already been moved
	if kingMoved {
		return castleMoves
	}

	// check if king side castling is possible if the rook has not been moved yet
	if !kingSideRookMoved {
		// [n,x,x,x,K,x,x,R] -> [x,x,x,x,x,R,K,x]
		kingSideEmpty := true
		for _, col := range []int{5, 6} {
			if g.board.grid[backrow][col] != nil {
				kingSideEmpty = false
			}
		}
		if kingSideEmpty {
			castleMoves = append(castleMoves, Position{backrow, 6})
		}
	}

	// check if queen side castling is possible if the rook has not been moved yet
	if !queenSideRookMoved {
		// [R,x,x,x,K,x,x,R] -> [x,x,K,R,x,x,x,x]
		queenSideEmpty := true
		for _, col := range []int{1, 2, 3} {
			if g.board.grid[backrow][col] != nil {
				queenSideEmpty = false
			}
		}
		if queenSideEmpty {
			castleMoves = append(castleMoves, Position{backrow, 2})
		}
	}

	return castleMoves
}

func (g *Game) getMovesForPawn(sourcePosition Position) []Position {
	var availableMoves []Position
	sourceRow, sourceCol := sourcePosition.Row, sourcePosition.Col
	piece := g.board.grid[sourceRow][sourceCol]
	if piece == nil {
		return availableMoves
	}

	currentPawnColor := piece.Color
	currentPawnDirection := 1
	if currentPawnColor == Black {
		currentPawnDirection = -1
	}

	// pawns can always move +1 or -1 forward based on the color of it
	forwardOne := Position{
		Row: sourceRow + currentPawnDirection,
		Col: sourceCol,
	}
	// if there's a piece directly in front of the pawn then it is blocked and cannot move forward.
	// only include the plus 1 forward if it's not blocked.
	if !g.board.IsPositionOutOfBounds(forwardOne) && g.board.grid[forwardOne.Row][forwardOne.Col] == nil {
		availableMoves = append(availableMoves, forwardOne)
	}

	// pawns can only move forward 2 spaces if it's the first turn.
	// check if the currentRow it is in is the initial pawn row for the color for this validation.
	// check if the rows leading up to it is blocked or not. if it's blocked then it can't traverse it as well
	if (piece.Color == Black && sourceRow == BlackPawnRow) ||
		(piece.Color == White && sourceRow == WhitePawnRow) {
		forwardTwo := Position{
			Row: sourceRow + (currentPawnDirection * 2),
			Col: sourceCol,
		}
		freePath := true
		for offset := 1; offset <= 2; offset++ {
			piece := g.board.grid[sourceRow+(currentPawnDirection*offset)][sourceCol]
			if piece != nil {
				freePath = false
			}
		}

		if freePath {
			availableMoves = append(availableMoves, forwardTwo)
		}
	}

	// pawns can move one space diagonally if there is a piece 1 space away from it on its left or right
	leftDiagonal := Position{
		Row: sourceRow + (currentPawnDirection),
		Col: sourceCol + -1,
	}
	if !g.board.IsPositionOutOfBounds(leftDiagonal) {
		leftDiagonalTarget := g.board.grid[leftDiagonal.Row][leftDiagonal.Col]
		if leftDiagonalTarget != nil && leftDiagonalTarget.Color != currentPawnColor {
			availableMoves = append(availableMoves, leftDiagonal)
		}
	}

	rightDiagonal := Position{
		Row: sourceRow + (currentPawnDirection),
		Col: sourceCol + 1,
	}
	if !g.board.IsPositionOutOfBounds(rightDiagonal) {
		rightDiagonalTarget := g.board.grid[rightDiagonal.Row][rightDiagonal.Col]
		if rightDiagonalTarget != nil && rightDiagonalTarget.Color != currentPawnColor {
			availableMoves = append(availableMoves, rightDiagonal)
		}
	}

	// en passant capture check, see: https://en.wikipedia.org/wiki/En_passant
	// this only handles the moves, it does not handle the opposing pawn removal if the player executes this move.
	// only check if there were moves made before in this game.
	moveCount := len(g.moves)
	if moveCount > 0 {
		lastMove := g.moves[moveCount-1]
		lastMovedPiece, lastMovedSourcePosition, lastMovedTargetPosition := lastMove.Piece, lastMove.SourcePosition, lastMove.TargetPosition

		// checking if the last moved piece was actually a pawn as en passant is only a valid move if it is performed directly after the oppose pawn's move.
		// unnecessary color check assuming the board state was updated correctly
		lastMovedPieceWasPawn := lastMovedPiece.PieceType == Pawn && lastMovedPiece.Color != currentPawnColor
		lastMovedPieceMovedTwoSpacesForward := int(math.Abs(float64(lastMovedTargetPosition.Row-lastMovedSourcePosition.Row))) == 2
		lastMovedPieceMovedStayedSameColumn := int(math.Abs(float64(lastMovedTargetPosition.Col-lastMovedSourcePosition.Col))) == 0

		// current pawn is on the same row as the opposing pawn and is adjacent to it on either side
		currentPawnAdjacentToOpposingPawn := sourceRow == lastMovedTargetPosition.Row && int(math.Abs(float64(sourceCol-lastMovedTargetPosition.Col))) == 1

		if lastMovedPieceWasPawn && lastMovedPieceMovedTwoSpacesForward && lastMovedPieceMovedStayedSameColumn && currentPawnAdjacentToOpposingPawn {
			enPassantPosition := Position{
				Row: sourceRow + currentPawnDirection,
				Col: lastMovedTargetPosition.Col,
			}
			availableMoves = append(availableMoves, enPassantPosition)
		}
	}

	return availableMoves
}

func (g *Game) getMovesForPiece(sourcePosition Position, moveSpec MoveSpec) []Position {
	var availableMoves []Position

	sourcePiece := g.board.grid[sourcePosition.Row][sourcePosition.Col]
	offsetsToCheck := moveSpec.Direction

	for _, offset := range offsetsToCheck {
		// for each offset direction, explore it until it's blocked by another piece or it's not a sliding piece which means it can only be moved once (King & Knight).
		// Pawns will have their own method and not use this as there are special rules for their movements
		nextPosition := sourcePosition
		for {
			nextPosition = Position{
				Row: nextPosition.Row + offset.DeltaRow,
				Col: nextPosition.Col + offset.DeltaCol,
			}
			if g.board.IsPositionOutOfBounds(nextPosition) {
				break
			}

			nextPositionPiece := g.board.grid[nextPosition.Row][nextPosition.Col]
			// if there's a piece blocking the next position, check if it's a possible move (opposing color piece, so the move is to remove it from the board)
			// if it's the same color, then it is not a possible move. either way, this direction offset exploration can be stopped
			if nextPositionPiece != nil {
				if nextPositionPiece.Color != sourcePiece.Color {
					availableMoves = append(availableMoves, nextPosition)
				}
				break
			}

			availableMoves = append(availableMoves, nextPosition)
			// Knight & King only
			if !moveSpec.Sliding {
				break
			}

		}
	}

	return availableMoves
}

func (g *Game) Print() {
	g.board.Print()
}
