package engine

import (
	"fmt"
	"strings"
)

type (
	Color     int8
	PieceType int8
)

const (
	A8, B8, C8, D8, E8, F8, G8, H8 = 57, 58, 59, 60, 61, 62, 63, 64
	A1, B1, C1, D1, E1, F1, G1, H1 = 1, 2, 3, 4, 5, 6, 7, 8
)

const (
	// Castling Right Mask
	WhiteKingsideRight  uint8 = 0x8
	WhiteQueensideRight uint8 = 0x4
	BlackKingsideRight  uint8 = 0x2
	BlackQueensideRight uint8 = 0x1
)

type State struct {
	EPFile         int8
	CastlingRights uint8
	Moved          Piece
	Captured       Piece
	Rule50         int8
	Ply            uint16
}

type Position struct {
	Board [65]Piece

	// Game state
	ColorToMove    Color
	CastlingRights uint8
	EPFile         int8
	Rule50         int8
	Ply            uint16

	// History
	prevStates  [100]State
	MoveHistory MoveList

	// King positions
	WhiteKing Square
	BlackKing Square

	// finished state
	Checkmate bool
	Stalemate bool

	// Piece Bitboards
	// in White pieces (P, N, B, R, Q, K)
	// 6-11 Black pieces (P, N, B, R, Q, K)
	pieceBitboards [2][6]Bitboard
}

func (pos *Position) PieceBitboard(p Piece) *Bitboard {
	return &pos.pieceBitboards[p.Color][p.PieceType]
}

func (pos *Position) ColorBitboard(c Color) Bitboard {
	switch c {
	case White:
		whiteBBs := pos.pieceBitboards[White]
		return whiteBBs[0] | whiteBBs[1] | whiteBBs[2] | whiteBBs[3] | whiteBBs[4] | whiteBBs[5]
	case Black:
		blackBBs := pos.pieceBitboards[Black]
		return blackBBs[0] | blackBBs[1] | blackBBs[2] | blackBBs[3] | blackBBs[4] | blackBBs[5]
	default:
		return 0
	}
}

func (pos *Position) AllPieces() Bitboard {
	return pos.ColorBitboard(White) | pos.ColorBitboard(Black)
}

// Call this after Board has been setup
func (pos *Position) InitPieceBitboards() {
	for i := Square(1); i <= 64; i++ {
		p := pos.Board[i]
		if p.PieceType != NoPiece {
			pos.pieceBitboards[p.Color][p.PieceType].Set(i)
		}
	}
}

func FromStandardStartingPosition() Position {
	return FromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
}

func (pos *Position) String() string {
	var str strings.Builder
	str.WriteString("---  Position --- \n")
	str.WriteString("Fen: " + FEN(pos) + "\n")

	for row := 7; row >= 0; row-- {
		str.WriteString("+---+---+---+---+---+---+---+---+\n")
		for file := 1; file <= 8; file++ {
			i := row*8 + file
			piece := pos.Board[i]
			str.WriteString(fmt.Sprintf("| %c ", PieceToFENChar[piece.Color][piece.PieceType]))
		}
		str.WriteString(fmt.Sprintf("| %d \n", row+1))
	}
	str.WriteString("+---+---+---+---+---+---+---+---+\n")
	str.WriteString("  a   b   c   d   e   f   g   h\n")
	str.WriteString("\n")

	return str.String()
}

// This function assumes that the move is valid
func (pos *Position) MakeMove(move Move) {
	movedPiece := pos.Board[move.From]
	capturedPiece := pos.Board[move.To]

	// Save the current state
	state := State{
		EPFile:         pos.EPFile,
		CastlingRights: pos.CastlingRights,
		Moved:          movedPiece,
		Captured:       capturedPiece,
		Rule50:         pos.Rule50,
		Ply:            pos.Ply,
	}

	// Move the piece
	pos.Board[move.To] = movedPiece
	pos.Board[move.From] = Piece{NoColor, NoPiece}

	// Update biboards
	fromBB := BBFromSquares(move.From)
	toBB := BBFromSquares(move.To)
	fromToBB := fromBB | toBB

	*pos.PieceBitboard(movedPiece) ^= fromToBB
	if capturedPiece.PieceType != NoPiece {
		*pos.PieceBitboard(capturedPiece) ^= toBB
	}

	// Update the EPFile

	pos.EPFile = 0

	// Rule 50
	pos.Rule50++
	if capturedPiece.PieceType != NoPiece || movedPiece.PieceType == Pawn {
		pos.Rule50 = 0
	}

	switch move.Flag {
	case PawnPush:
		pos.EPFile = File(move.From)
	case Castling:
		if movedPiece.Color == White {
			if move.To == Square(3) {
				pos.swapSquares(Square(1), Square(4))
			}
			if move.To == Square(7) {
				pos.swapSquares(Square(8), Square(6))
			}
		}
		if movedPiece.Color == Black {
			if move.To == Square(59) {
				pos.swapSquares(Square(57), Square(60))
			}
			if move.To == Square(63) {
				pos.swapSquares(Square(64), Square(62))
			}
		}
	case EnPassentCapture:
		if movedPiece.Color == White {
			pos.Board[move.To-8] = Piece{NoColor, NoPiece}
			*pos.PieceBitboard(Piece{Black, Pawn}) ^= BBFromSquares(move.To - 8)
			state.Captured = Piece{PieceType: Pawn, Color: Black}
		} else {
			pos.Board[move.To+8] = Piece{NoColor, NoPiece}
			*pos.PieceBitboard(Piece{White, Pawn}) ^= BBFromSquares(move.To + 8)
			state.Captured = Piece{PieceType: Pawn, Color: White}
		}
	case PromotionToQueen:
		pos.Board[move.To] = Piece{Color: movedPiece.Color, PieceType: Queen}
		*pos.PieceBitboard(Piece{movedPiece.Color, Queen}) ^= toBB
		*pos.PieceBitboard(movedPiece) ^= toBB
	case PromotionToKnight:
		pos.Board[move.To] = Piece{Color: movedPiece.Color, PieceType: Knight}
		*pos.PieceBitboard(Piece{movedPiece.Color, Knight}) ^= toBB
		*pos.PieceBitboard(movedPiece) ^= toBB
	case PromotionToRook:
		pos.Board[move.To] = Piece{Color: movedPiece.Color, PieceType: Rook}
		*pos.PieceBitboard(Piece{movedPiece.Color, Rook}) ^= toBB
		*pos.PieceBitboard(movedPiece) ^= toBB
	case PromotionToBishop:
		pos.Board[move.To] = Piece{Color: movedPiece.Color, PieceType: Bishop}
		*pos.PieceBitboard(Piece{movedPiece.Color, Bishop}) ^= toBB
		*pos.PieceBitboard(movedPiece) ^= toBB
	}

	pos.updateCastlingRights()

	pos.prevStates[pos.Ply] = state
	pos.Ply++

	if movedPiece.PieceType == King {
		if movedPiece.Color == White {
			pos.WhiteKing = move.To
		} else {
			pos.BlackKing = move.To
		}
	}

	pos.ColorToMove = pos.ColorToMove.opposite()
}

func (pos *Position) updateCastlingRights() {
	wk, wq, bk, bq := pos.getCastlingRights()
	// Whitekingside
	if wk &&
		(pos.Board[5].Equal(Piece{White, King})) &&
		(pos.Board[8].Equal(Piece{White, Rook})) {
		pos.CastlingRights |= WhiteKingsideRight
	} else {
		pos.CastlingRights &= ^WhiteKingsideRight
	}

	// WhiteQueenside
	if wq &&
		(pos.Board[5].Equal(Piece{White, King})) &&
		(pos.Board[1].Equal(Piece{White, Rook})) {
		pos.CastlingRights |= WhiteQueensideRight
	} else {
		pos.CastlingRights &= ^WhiteQueensideRight
	}

	// Blackkingside
	if bk &&
		(pos.Board[61].Equal(Piece{Black, King})) &&
		(pos.Board[64].Equal(Piece{Black, Rook})) {
		pos.CastlingRights |= BlackKingsideRight
	} else {
		pos.CastlingRights &= ^BlackKingsideRight
	}

	// BlackQueenside
	if bq &&
		(pos.Board[61].Equal(Piece{Black, King})) &&
		(pos.Board[57].Equal(Piece{Black, Rook})) {
		pos.CastlingRights |= BlackQueensideRight
	} else {
		pos.CastlingRights &= ^BlackQueensideRight
	}
}

func (pos *Position) UndoMove(move Move) {
	pos.Ply--
	prevState := pos.prevStates[pos.Ply]

	pos.EPFile = prevState.EPFile
	pos.Rule50 = prevState.Rule50
	pos.CastlingRights = prevState.CastlingRights

	pos.Board[move.From] = prevState.Moved

	// Update biboards
	fromBB := BBFromSquares(move.From)
	toBB := BBFromSquares(move.To)
	fromToBB := fromBB | toBB

	*pos.PieceBitboard(prevState.Moved) ^= fromToBB

	switch move.Flag {
	default:
		pos.Board[move.To] = prevState.Captured
		if prevState.Captured.PieceType != NoPiece {
			*pos.PieceBitboard(prevState.Captured) ^= toBB
		}
	case Castling:
		if prevState.Moved.Color == White {
			if move.To == Square(3) {
				pos.swapSquares(Square(1), Square(4))
				pos.Board[3] = ANoPiece()
				*pos.PieceBitboard(Piece{White, King}) ^= BBFromSquares(3, 4)
			}
			if move.To == Square(7) {
				pos.swapSquares(Square(8), Square(6))
				pos.Board[7] = ANoPiece()
				*pos.PieceBitboard(Piece{White, King}) ^= BBFromSquares(5, 6)
			}
		} else {
			if move.To == Square(59) {
				pos.swapSquares(Square(57), Square(60))
				pos.Board[59] = ANoPiece()
				*pos.PieceBitboard(Piece{Black, King}) ^= BBFromSquares(59, 60)
			}
			if move.To == Square(63) {
				pos.swapSquares(Square(64), Square(62))
				pos.Board[63] = ANoPiece()
				*pos.PieceBitboard(Piece{Black, King}) ^= BBFromSquares(61, 62)
			}
		}
	case EnPassentCapture:
		pos.Board[move.To] = Piece{NoColor, NoPiece}
		if prevState.Moved.Color == White {
			pos.Board[move.To-8] = Piece{Black, Pawn}
			*pos.PieceBitboard(Piece{Black, Pawn}) ^= BBFromSquares(move.To - 8)
		} else {
			pos.Board[move.To+8] = Piece{White, Pawn}
			*pos.PieceBitboard(Piece{White, Pawn}) ^= BBFromSquares(move.To + 8)
		}
	}

	if prevState.Moved.PieceType == King {
		if prevState.Moved.Color == White {
			pos.WhiteKing = move.From
		} else {
			pos.BlackKing = move.From
		}
	}

	pos.ColorToMove = pos.ColorToMove.opposite()
}

// order : wk, wq, bk, bq
func (pos *Position) getCastlingRights() (bool, bool, bool, bool) {
	return (pos.CastlingRights&WhiteKingsideRight != 0),
		(pos.CastlingRights&WhiteQueensideRight != 0),
		(pos.CastlingRights&BlackKingsideRight != 0),
		(pos.CastlingRights&BlackQueensideRight != 0)
}

func (pos *Position) GetKingSquare(color Color) Square {
	if color == White {
		return pos.WhiteKing
	}
	return pos.BlackKing
}

func (pos *Position) FlipColor() {
	pos.ColorToMove = pos.ColorToMove.opposite()
}

func (pos *Position) swapSquares(a, b Square) {
	pos.Board[a], pos.Board[b] = pos.Board[b], pos.Board[a]

	// Update bitboards, disgusting
	aBB := BBFromSquares(a)
	bBB := BBFromSquares(b)
	abBB := aBB | bBB
	aPiece := pos.Board[a]
	bPiece := pos.Board[b]
	if aPiece.PieceType != NoPiece {
		*pos.PieceBitboard(aPiece) ^= abBB
	}
	if bPiece.PieceType != NoPiece {
		*pos.PieceBitboard(bPiece) ^= abBB
	}
}
