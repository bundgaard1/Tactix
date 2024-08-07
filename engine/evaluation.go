package engine

const (
	// Piece Values
	PawnValue   = 100
	KnightValue = 300
	BishopValue = 330
	RookValue   = 500
	QueenValue  = 900

	CheckmateValue = 10000
)

// positive for white, negative for black, as it should be
func Evaluate(pos *Position) int {
	eval := 0

	if pos.Checkmate {
		eval = CheckmateValue
	}

	eval *= who2move(pos.ColorToMove)

	return eval
}

func who2move(c2m Color) int {
	if c2m == White {
		return 1
	}
	return -1
}
