package engine

import "fmt"

const (
	// Piece types
	Pawn    PieceType = 0
	Knight  PieceType = 1
	Bishop  PieceType = 2
	Rook    PieceType = 3
	Queen   PieceType = 4
	King    PieceType = 5
	NoPiece PieceType = 6

	// Colors
	White   Color = 0
	Black   Color = 1
	NoColor Color = 2
)

type Piece struct {
	Color
	PieceType
}

func (piece *Piece) Unwrap() (Color, PieceType) {
	return piece.Color, piece.PieceType
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

func WhitePiece(p PieceType) Piece {
	return Piece{White, p}
}

func BlackPiece(p PieceType) Piece {
	return Piece{Black, p}
}

func ANoPiece() Piece {
	return Piece{NoColor, NoPiece}
}
