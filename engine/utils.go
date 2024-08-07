package engine

import (
	"errors"
	"fmt"
	"strings"
)

func DeriveSquare(file int, rank int) Square {
	return Square((rank-1)*8 + file)
}

func Rank(square Square) int8 {
	return int8(square-1)/8 + 1
}

var FileRune = map[int8]rune{
	1: 'a',
	2: 'b',
	3: 'c',
	4: 'd',
	5: 'e',
	6: 'f',
	7: 'g',
	8: 'h',
}

func File(square Square) int8 {
	return int8(square-1)%8 + 1
}

func CastelingRightsToString(rights uint8) string {
	castingString := ""
	if rights == 0 {
		castingString = "-"
	} else {
		if rights&WhiteKingsideRight != 0 {
			castingString += "K"
		}
		if rights&WhiteQueensideRight != 0 {
			castingString += "Q"
		}
		if rights&BlackKingsideRight != 0 {
			castingString += "k"
		}
		if rights&BlackQueensideRight != 0 {
			castingString += "q"
		}
	}
	return castingString
}

var ErrInvalidMove = errors.New("invalid move")

// ParseMove in the form like "e2e4" to a Move struct
func ParseUCIMove(pos *Position, uciMove string) (Move, error) {
	if len(uciMove) < 4 || 5 < len(uciMove) {
		return Move{}, ErrInvalidMove
	}

	var move Move
	from := uciMove[0:2]
	to := uciMove[2:4]

	fromFile := from[0] - 'a' + 1
	move.From = Square((from[1]-'0'-1)*8 + fromFile)

	toFile := to[0] - 'a' + 1
	move.To = Square((to[1]-'0'-1)*8 + toFile)

	if len(uciMove) == 5 {
		switch uciMove[5] {
		case 'q':
			move.Flag = PromotionToQueen
		case 'r':
			move.Flag = PromotionToRook
		case 'b':
			move.Flag = PromotionToBishop
		case 'n':
			move.Flag = PromotionToKnight
		default:
			return Move{}, ErrInvalidMove
		}
	} else {
		move.Flag = flagForMove(pos, move)
	}

	return move, nil
}

// Promotions are handled
func flagForMove(pos *Position, move Move) MoveFlag {
	if pos.Board[move.From].PieceType == Pawn {
		if pos.Board[move.To].PieceType == NoPiece {
			if File(move.To) == pos.EPFile {
				return EnPassentCapture
			}
			if Rank(move.To) == Rank(move.From)+2 || Rank(move.To) == Rank(move.From)-2 {
				return PawnPush
			}
		}
		if Rank(move.To) == 1 || Rank(move.To) == 8 {
			return PromotionToQueen
		}
	} else if pos.Board[move.From].PieceType == King {
		if move.From == E1 && move.To == G1 {
			return Castling
		}
		if move.From == E1 && move.To == C1 {
			return Castling
		}
		if move.From == E8 && move.To == G8 {
			return Castling
		}
		if move.From == E8 && move.To == C8 {
			return Castling
		}
	}

	return NoFlag
}

func IsMoveValid(pos *Position, move Move) bool {
	if move.From == 0 || move.To == 0 {
		return false
	}

	if pos.Board[move.From].PieceType == NoPiece {
		return false
	}

	if pos.Board[move.From].Color != pos.ColorToMove {
		return false
	}

	if pos.Board[move.To].Color == pos.ColorToMove {
		return false
	}

	// Check if the move is valid
	moveList := LegalMoves(pos)
	for i := 0; i < moveList.Count; i++ {
		if moveList.Moves[i] == move {
			return true
		}
	}
	return false
}

func (sq Square) String() string {
	return fmt.Sprintf("%c%d", FileRune[File(sq)], Rank(sq))
}

func (m Move) String() string {
	return fmt.Sprintf("From: %d, To: %d, Flag: %d", m.From, m.To, m.Flag)
}

func (m Move) UCIString() string {
	var strbuilder strings.Builder

	strbuilder.WriteString(fmt.Sprintf("%c%d%c%d", FileRune[File(m.From)], Rank(m.From), FileRune[File(m.To)], Rank(m.To)))
	switch m.Flag {
	case PromotionToQueen:
		strbuilder.WriteString("q")
	case PromotionToBishop:
		strbuilder.WriteString("b")
	case PromotionToKnight:
		strbuilder.WriteString("n")
	case PromotionToRook:
		strbuilder.WriteString("r")
	}
	return strbuilder.String()
}
