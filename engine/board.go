package engine

import "fmt"

type Color uint8
type PieceType uint8

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

var state struct {
}

// One extra for invalid: index 0
// 1-8 rank 1
// ...
// 56-64 rank 8
var Position [65]Piece

func FromStandardStartingPosition() {
	Position[1], Position[2], Position[3], Position[4] = Piece{White, Rook}, Piece{White, Knight}, Piece{White, Bishop}, Piece{White, Queen}
	Position[5], Position[6], Position[7], Position[8] = Piece{White, King}, Piece{White, Bishop}, Piece{White, Knight}, Piece{White, Rook}
	Position[9], Position[10], Position[11], Position[12], Position[13], Position[14], Position[15], Position[16] = Piece{White, Pawn}, Piece{White, Pawn}, Piece{White, Pawn}, Piece{White, Pawn}, Piece{White, Pawn}, Piece{White, Pawn}, Piece{White, Pawn}, Piece{White, Pawn}

	Position[57], Position[58], Position[59], Position[60] = Piece{Black, Rook}, Piece{Black, Knight}, Piece{Black, Bishop}, Piece{Black, Queen}
	Position[61], Position[62], Position[63], Position[64] = Piece{Black, King}, Piece{Black, Bishop}, Piece{Black, Knight}, Piece{Black, Rook}
	Position[49], Position[50], Position[51], Position[52], Position[53], Position[54], Position[55], Position[56] = Piece{Black, Pawn}, Piece{Black, Pawn}, Piece{Black, Pawn}, Piece{Black, Pawn}, Piece{Black, Pawn}, Piece{Black, Pawn}, Piece{Black, Pawn}, Piece{Black, Pawn}
}

func PrintPosition() {
	for row := 7; row >= 0; row-- {
		for file := 1; file <= 8; file++ {
			fmt.Print(Position[row*8+file])
		}
		fmt.Print("\n")
	}
}
