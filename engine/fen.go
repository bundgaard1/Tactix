package engine

import (
	"fmt"
	"strings"
)

const StartingPositionFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

var FENCharToPiece = map[rune]Piece{
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

var PieceToFENChar = map[Color]map[PieceType]rune{
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
		King:   'k',
		Queen:  'q',
	},
}

// Should comply with FEN standard
func FromFEN(fen string) Position {
	var pos Position
	fenFields := strings.Split(fen, " ")

	if len(fenFields) != 6 {
		panic("Not 6 fields in FEN-string, got " + fmt.Sprint(len(fenFields)))
	}

	rank := 8
	file := 1

	for _, char := range fenFields[0] {
		switch char {
		case '/':
			rank--
			file = 1
		case 'P', 'N', 'B', 'R', 'Q', 'K', 'p', 'n', 'b', 'r', 'q', 'k':
			piece := FENCharToPiece[char]
			sq := DeriveSquare(file, rank)
			pos.Board[sq] = piece
			file++
			if char == 'K' {
				pos.WhiteKing = sq
			} else if char == 'k' {
				pos.BlackKing = sq
			}
		case '1', '2', '3', '4', '5', '6', '7', '8':
			file += int(char - '0')
		default:
			panic("Invalid FEN: board position.")
		}

	}
	// Check if both kings are present
	if pos.WhiteKing == 0 || pos.BlackKing == 0 {
		panic("Invalid FEN: missing king.")
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

	pos.checkmate = false
	pos.stalemate = false

	return pos
}

// Should comply with FEN standard
func FEN(pos *Position) string {
	var fen strings.Builder

	// Board Position
	for rank := 8; rank >= 1; rank-- {
		empty := 0
		for file := 1; file <= 8; file++ {
			piece := pos.Board[(rank-1)*8+file]
			if piece.PieceType == NoPiece {
				empty++
			} else {
				if empty > 0 {
					fen.WriteString(fmt.Sprint(empty))
					empty = 0
				}
				fen.WriteRune(PieceToFENChar[piece.Color][piece.PieceType])
			}
		}
		if empty > 0 {
			fen.WriteString(fmt.Sprint(empty))
		}
		if rank > 1 {
			fen.WriteRune('/')
		}
	}
	fen.WriteRune(' ')

	// Side to move
	if pos.ColorToMove == White {
		fen.WriteRune('w')
	} else {
		fen.WriteRune('b')
	}
	fen.WriteRune(' ')

	// Castling rights
	if pos.CastlingRights == 0 {
		fen.WriteRune('-')
	} else {
		wk, wq, bk, bq := pos.getCastlingRights()
		if wk {
			fen.WriteRune('K')
		}
		if wq {
			fen.WriteRune('Q')
		}
		if bk {
			fen.WriteRune('k')
		}
		if bq {
			fen.WriteRune('q')
		}
	}
	fen.WriteRune(' ')

	// En passant target square
	if pos.EPFile == 0 {
		fen.WriteRune('-')
	} else {
		fen.WriteRune('a' + rune(pos.EPFile))
		if pos.ColorToMove == White {
			fen.WriteRune('3')
		} else {
			fen.WriteRune('6')
		}
	}

	fen.WriteRune(' ')

	// Halfmove clock
	fen.WriteString(fmt.Sprint(pos.Rule50))
	fen.WriteRune(' ')

	// Fullmove counter
	fen.WriteString(fmt.Sprint(pos.Ply))

	return fen.String()
}
