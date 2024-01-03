package engine

import "fmt"

type Color int8
type PieceType int8

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

// One extra for invalid: index 0
// 1-8 rank 1
// ...
// 56-64 rank 8
type Position struct {
	Board       [65]Piece
	ColorToMove Color

	// Castling allowance
	BlackKingCastle  bool
	BlackQueenCastle bool
	WhiteKingCastle  bool
	WhiteQueenCastle bool
}

func (pos *Position) FromStandardStartingPosition() {
	pos.ColorToMove = White

	pos.Board[1], pos.Board[2], pos.Board[3], pos.Board[4] = Piece{White, Rook}, Piece{White, Knight}, Piece{White, Bishop}, Piece{White, Queen}
	pos.Board[5], pos.Board[6], pos.Board[7], pos.Board[8] = Piece{White, King}, Piece{White, Bishop}, Piece{White, Knight}, Piece{White, Rook}
	for i := 9; i <= 16; i++ {
		pos.Board[i] = Piece{White, Pawn}
	}

	for i := 49; i <= 56; i++ {
		pos.Board[i] = Piece{Black, Pawn}
	}

	pos.Board[57], pos.Board[58], pos.Board[59], pos.Board[60] = Piece{Black, Rook}, Piece{Black, Knight}, Piece{Black, Bishop}, Piece{Black, Queen}
	pos.Board[61], pos.Board[62], pos.Board[63], pos.Board[64] = Piece{Black, King}, Piece{Black, Bishop}, Piece{Black, Knight}, Piece{Black, Rook}
}

func (pos *Position) PrintPosition() {
	for row := 7; row >= 0; row-- {
		for file := 1; file <= 8; file++ {
			i := row*8 + file
			fmt.Printf(" %02d:", i)
			fmt.Print(pos.Board[i])
		}
		fmt.Print("\n")
	}
}

func Rank(square int8) int8 {
	return (square-1)/8 + 1
}

func File(square int8) int8 {
	return (square-1)%8 + 1
}
