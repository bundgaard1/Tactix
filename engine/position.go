package engine

import (
	"fmt"
	"strings"
)

type Color int8
type PieceType int8

const (
	A8, B8, C8, D8, E8, F8, G8, H8 = 57, 58, 59, 60, 61, 62, 63, 64
	A7, B7, C7, D7, E7, F7, G7, H7 = 49, 50, 51, 52, 53, 54, 55, 56
	A6, B6, C6, D6, E6, F6, G6, H6 = 41, 42, 43, 44, 45, 46, 47, 48
	A5, B5, C5, D5, E5, F5, G5, H5 = 33, 34, 35, 36, 37, 38, 39, 40
	A4, B4, C4, D4, E4, F4, G4, H4 = 25, 26, 27, 28, 29, 30, 31, 32
	A3, B3, C3, D3, E3, F3, G3, H3 = 17, 18, 19, 20, 21, 22, 23, 24
	A2, B2, C2, D2, E2, F2, G2, H2 = 9, 10, 11, 12, 13, 14, 15, 16
	A1, B1, C1, D1, E1, F1, G1, H1 = 1, 2, 3, 4, 5, 6, 7, 8
)

const (
	// Piece types
	NoPiece PieceType = 0
	Pawn    PieceType = 1
	Knight  PieceType = 2
	Bishop  PieceType = 3
	Rook    PieceType = 4
	King    PieceType = 5
	Queen   PieceType = 6

	// Colors
	NoColor Color = 0
	White   Color = 1
	Black   Color = 2

	// Castling Right Mask
	WhiteKingsideRight  uint8 = 0x8
	WhiteQueensideRight uint8 = 0x4
	BlackKingsideRight  uint8 = 0x2
	BlackQueensideRight uint8 = 0x1
)

func (c Color) String() string {
	if c == White {
		return "White"
	} else if c == Black {
		return "Black"
	}
	return "NO COLOR"
}

func (c Color) opposite() Color {
	if c == White {
		return Black
	} else if c == Black {
		return White
	}
	return NoColor
}

type Piece struct {
	Color
	PieceType
}

type State struct {
	EPFile         int8
	CastlingRights uint8
	Moved          Piece
	Captured       Piece
	StatePly       uint16
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
	prevStates [100]State

	// Kings position
	WhiteKing Square
	BlackKing Square
}

func FromStandardStartingPosition() Position {
	return FromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
}

func (pos *Position) String() string {
	var str strings.Builder
	str.WriteString("---  Position --- \n")
	str.WriteString(fmt.Sprintf("	CTM: %s \n	EP-file: %d \n 	Ply: %d \n 	50 Rule: %d \n", pos.ColorToMove, pos.EPFile, pos.Ply, pos.Rule50))
	str.WriteString(fmt.Sprintf(" 	Castling: %s \n", CastelingRightsToString(pos.CastlingRights)))

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
	capturedPiece := pos.Board[move.To] // Or no piece

	// Save the current state
	state := State{
		EPFile:         pos.EPFile,
		CastlingRights: pos.CastlingRights,
		Moved:          movedPiece,
		Captured:       capturedPiece,
		StatePly:       pos.Ply,
	}

	// Move the piece
	pos.Board[move.To] = movedPiece
	pos.Board[move.From] = Piece{NoColor, NoPiece}

	pos.EPFile = 0

	switch move.Flag {
	case PawnPush:
		pos.EPFile = File(move.From)
	case Castling:
		if movedPiece.Color == White {
			if move.To == Square(3) {
				pos.Board[4] = pos.Board[1]
				pos.Board[1] = Piece{NoColor, NoPiece}
			}
			if move.To == Square(7) {
				pos.Board[6] = pos.Board[8]
				pos.Board[8] = Piece{NoColor, NoPiece}
			}
		}
		if movedPiece.Color == Black {
			if move.To == Square(59) {
				pos.Board[60] = pos.Board[57]
				pos.Board[57] = Piece{NoColor, NoPiece}
			}
			if move.To == Square(63) {
				pos.Board[62] = pos.Board[64]
				pos.Board[64] = Piece{NoColor, NoPiece}
			}
		}
	case EnPassentCapture:
		if movedPiece.Color == White {
			pos.Board[move.To-8] = Piece{NoColor, NoPiece}
			state.Captured = Piece{PieceType: Pawn, Color: Black}
		} else {
			pos.Board[move.To+8] = Piece{NoColor, NoPiece}
			state.Captured = Piece{PieceType: Pawn, Color: White}
		}
	case PromotionToQueen:
		pos.Board[move.To] = Piece{Color: movedPiece.Color, PieceType: Queen}
	case PromotionToKnight:
		pos.Board[move.To] = Piece{Color: movedPiece.Color, PieceType: Knight}
	case PromotionToRook:
		pos.Board[move.To] = Piece{Color: movedPiece.Color, PieceType: Rook}
	case PromotionToBishop:
		pos.Board[move.To] = Piece{Color: movedPiece.Color, PieceType: Bishop}

	}

	// Moving rook or king removes rights
	if movedPiece.PieceType == King {
		if movedPiece.Color == White {
			pos.CastlingRights &= 0b0011
		}
		if movedPiece.Color == Black {
			pos.CastlingRights &= 0b1100
		}
	}
	if movedPiece.PieceType == Rook {
		if movedPiece.Color == White {
			if move.From == 1 {
				pos.CastlingRights &= 0b1011
			} else if move.From == 8 {
				pos.CastlingRights &= 0b0111
			}
		}
		if movedPiece.Color == Black {
			if move.From == 57 {
				pos.CastlingRights &= 0b1110
			} else if move.From == 64 {
				pos.CastlingRights &= 0b1101
			}
		}
	}

	// Update the king position
	if movedPiece.PieceType == King {
		pos.updateKingSquare(movedPiece.Color, Square(move.To))
	}

	pos.prevStates[pos.Ply] = state
	pos.Ply++

	pos.ColorToMove = pos.ColorToMove.opposite()

}

func (pos *Position) UndoMove(move Move) {
	pos.Ply--
	prevState := pos.prevStates[pos.Ply]

	pos.EPFile = prevState.EPFile
	pos.CastlingRights = prevState.CastlingRights

	pos.Board[move.From] = prevState.Moved

	switch move.Flag {
	default:
		pos.Board[move.To] = prevState.Captured
	case Castling:
		if prevState.Moved.Color == White {
			if move.To == Square(3) {
				pos.Board[1] = pos.Board[4]
				pos.Board[3] = Piece{NoColor, NoPiece}
				pos.Board[4] = Piece{NoColor, NoPiece}
			}
			if move.To == Square(7) {
				pos.Board[8] = pos.Board[6]
				pos.Board[6] = Piece{NoColor, NoPiece}
				pos.Board[7] = Piece{NoColor, NoPiece}
			}
		} else {
			if move.To == Square(59) {
				pos.Board[57] = pos.Board[60]
				pos.Board[60] = Piece{NoColor, NoPiece}
				pos.Board[59] = Piece{NoColor, NoPiece}
			}
			if move.To == Square(63) {
				pos.Board[64] = pos.Board[62]
				pos.Board[62] = Piece{NoColor, NoPiece}
				pos.Board[63] = Piece{NoColor, NoPiece}
			}
		}
	case EnPassentCapture:
		pos.Board[move.To] = Piece{NoColor, NoPiece}
		if prevState.Moved.Color == White {
			pos.Board[move.To-8] = Piece{Black, Pawn}
		} else {
			pos.Board[move.To+8] = Piece{White, Pawn}
		}
	}

	// Update the king position
	movedPiece := prevState.Moved
	if movedPiece.PieceType == King {
		pos.updateKingSquare(movedPiece.Color, Square(move.From))
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

func (pos *Position) updateKingSquare(Color Color, square Square) {
	if Color == White {
		pos.WhiteKing = square
	} else {
		pos.BlackKing = square
	}

}
