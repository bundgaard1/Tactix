package engine

import (
	"fmt"
	"strings"
)

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
	// One extra for invalid: index 0
	// 1-8 rank 1
	// ...
	// 56-64 rank 8
	Board          [65]Piece
	ColorToMove    Color
	CastlingRights uint8
	EPFile         int8
	Rule50         int8
	Ply            uint16
	prevStates     [100]State
}

var CharToPiece = map[rune]Piece{
	'P': {White, Pawn},
	'N': {White, Knight},
	'B': {White, Bishop},
	'R': {White, Rook},
	'Q': {White, Queen},
	'K': {White, King},
	'p': {Black, Pawn},
	'n': {Black, Knight},
	'b': {Black, Bishop},
	'r': {Black, Rook},
	'q': {Black, Queen},
	'k': {Black, King},
}

func (pos *Position) FromFEN(fen string) {
	fenFields := strings.Split(fen, " ")

	if len(fenFields) != 6 {
		panic("Not 6 fields in FEN-string")
	}

	rank := 8
	file := 1

	for _, char := range fenFields[0] {
		switch char {
		case '/':
			rank--
			file = 1
		case 'P', 'N', 'B', 'R', 'Q', 'K', 'p', 'n', 'b', 'r', 'q', 'k':
			piece := CharToPiece[char]
			pos.Board[(rank-1)*8+file] = piece
			file++
		case '1', '2', '3', '4', '5', '6', '7', '8':
			file += int(char - '0')
		}

	}

	// Side to move

	if fenFields[1] == "w" {
		pos.ColorToMove = White
	} else if fenFields[1] == "b" {
		pos.ColorToMove = Black
	} else {
		panic("Invalid FEN: side to move.")
	}

	// Castling ability

	pos.CastlingRights = 0

	for _, char := range fenFields[2] {
		switch char {
		case 'K':
			pos.CastlingRights |= WhiteKingsideRight
		case 'Q':
			pos.CastlingRights |= WhiteQueensideRight
		case 'k':
			pos.CastlingRights |= BlackKingsideRight
		case 'q':
			pos.CastlingRights |= BlackQueensideRight
		}
	}

	// En passant target square

	if fenFields[3][0] == '-' {
		pos.EPFile = 0
	} else {
		pos.EPFile = int8(fenFields[3][0] - '0')
	}

	// halfmove clock
	pos.Rule50 = int8(fenFields[4][0] - '0')

	// Fullmove counter
	pos.Ply = uint16(fenFields[5][0] - '0')
}

func (pos *Position) FromStandardStartingPosition() {
	pos.ColorToMove = White
	pos.CastlingRights = 0b1111
	pos.EPFile = 0
	pos.Ply = 1

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

var PieceToChar = map[Color]map[PieceType]rune{
	NoColor: {
		NoPiece: ' ',
	},
	White: {
		Pawn:   'P',
		Knight: 'N',
		Bishop: 'B',
		Rook:   'R',
		King:   'K',
		Queen:  'Q',
	},
	Black: {
		Pawn:   'p',
		Knight: 'n',
		Bishop: 'b',
		Rook:   'r',
		King:   'K',
		Queen:  'q',
	},
}

func (pos *Position) PrintPosition() {
	fmt.Printf("---  Position --- \n")
	fmt.Printf("	CTM: %s \n	EP-file: %d \n 	Ply: %d \n 	50 Rule: %d \n", pos.ColorToMove, pos.EPFile, pos.Ply, pos.Rule50)
	fmt.Printf(" 	Castling: %s \n", CastelingRightsToString(pos.CastlingRights))

	for row := 7; row >= 0; row-- {
		fmt.Printf("+---+---+---+---+---+---+---+---+\n")
		for file := 1; file <= 8; file++ {
			i := row*8 + file
			piece := pos.Board[i]
			fmt.Printf("| %c ", PieceToChar[piece.Color][piece.PieceType])
		}
		fmt.Printf("| %d \n", row+1)
	}
	fmt.Printf("+---+---+---+---+---+---+---+---+\n")
	fmt.Printf("  a   b   c   d   e   f   g   h\n")
	fmt.Printf("\n")
}

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

	pos.ColorToMove = pos.ColorToMove.opposite()
}
