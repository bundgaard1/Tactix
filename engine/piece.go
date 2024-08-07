package engine

import "fmt"

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
)

type Piece struct {
	Color
	PieceType
}

func (p Piece) String() string {
	return fmt.Sprintf("%s%s", p.Color.String(), p.PieceType.String())
}

func (p Piece) Equal(other Piece) bool {
	return p.Color == other.Color && p.PieceType == other.PieceType
}

func (p PieceType) String() string {
	switch p {
	case Pawn:
		return "P"
	case Knight:
		return "N"
	case Bishop:
		return "B"
	case Rook:
		return "R"
	case King:
		return "K"
	case Queen:
		return "Q"
	}
	return "NO PIECE"
}

func (c Color) String() string {
	if c == White {
		return "w"
	} else if c == Black {
		return "b"
	}
	return "x"
}

func (c Color) opposite() Color {
	switch c {
	case White:
		return Black
	case Black:
		return White
	default:
		return NoColor
	}
}

func WhitePiece(p PieceType) Piece {
	return Piece{White, p}
}

func BlackPiece(p PieceType) Piece {
	return Piece{Black, p}
}

func ANoPiece() Piece {
	return Piece{NoColor, NoPiece}
}
