package engine

import "fmt"

type (
	Color int8
	PType int8
)

const (
	// Piece types
	Pawn    PType = 0
	Knight  PType = 1
	Bishop  PType = 2
	Rook    PType = 3
	Queen   PType = 4
	King    PType = 5
	NoPiece PType = 6

	// Colors
	White   Color = 0
	Black   Color = 1
	NoColor Color = 2
)

type Piece struct {
	Color
	PType
}

func (piece *Piece) Unwrap() (Color, PType) {
	return piece.Color, piece.PType
}

func (p Piece) String() string {
	return fmt.Sprintf("%s%s", p.Color.String(), p.PType.String())
}

func (p Piece) Equal(other Piece) bool {
	return p.Color == other.Color && p.PType == other.PType
}

func (p PType) String() string {
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
	switch c {
	case White:
		return "w"
	case Black:
		return "b"
	default:
		return "x"
	}
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

func WhitePiece(p PType) Piece {
	return Piece{White, p}
}

func BlackPiece(p PType) Piece {
	return Piece{Black, p}
}

func ANoPiece() Piece {
	return Piece{NoColor, NoPiece}
}
