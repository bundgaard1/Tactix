package engine

import "errors"

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
// Promotion is not supported (yet)
func ParseUCIMove(pos *Position, algMove string) (Move, error) {
	if len(algMove) != 4 {
		return Move{}, ErrInvalidMove
	}

	move := Move{}
	from := algMove[:2]
	to := algMove[2:]

	fromFile := from[0] - 'a' + 1
	move.From = Square((from[1]-'0'-1)*8 + fromFile)

	toFile := to[0] - 'a' + 1
	move.To = Square((to[1]-'0'-1)*8 + toFile)

	move.Flag = flagForMove(pos, move)

	return move, nil
}

// does not support different promotion
func flagForMove(pos *Position, move Move) MoveFlag {
	if pos.Board[move.From].PieceType == Pawn {
		if pos.Board[move.To].PieceType == NoPiece {
			if File(move.To) == pos.EPFile {
				return EnPassentCapture
			}
			return PawnPush
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
	} else if pos.Board[move.From].PieceType == Pawn && (Rank(move.To) == 1 || Rank(move.To) == 8) {
		return PromotionToQueen
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
	moveList := GetValidMoves(pos)
	for i := 0; i < moveList.Count; i++ {
		if moveList.Moves[i] == move {
			return true
		}
	}
	return false
}
