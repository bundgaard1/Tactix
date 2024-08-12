package engine

const (
	// Piece Values
	PawnValue   = 100
	KnightValue = 300
	BishopValue = 320
	RookValue   = 500
	QueenValue  = 900

	MobilityValue = 10

	PositiveInfinity = 999_999
	NegativeInfinity = -PositiveInfinity
)

// positive for white, negative for black, as it should be
func Evaluate(pos *Position) int {
	eval := 0

	if pos.Checkmate {
		return PositiveInfinity * who2move(pos.ColorToMove)
	}

	materialScore := materialScore(pos)
	mobilityScore := mobilityScore(pos)

	eval = (materialScore + mobilityScore)

	return eval
}

func materialScore(pos *Position) (materialScore int) {
	materialScore = PawnValue * (pos.PieceBitboard(WhitePiece(Pawn)).Count() - pos.PieceBitboard(BlackPiece(Pawn)).Count())
	materialScore += KnightValue * (pos.PieceBitboard(WhitePiece(Knight)).Count() - pos.PieceBitboard(BlackPiece(Knight)).Count())
	materialScore += BishopValue * (pos.PieceBitboard(WhitePiece(Bishop)).Count() - pos.PieceBitboard(BlackPiece(Bishop)).Count())
	materialScore += RookValue * (pos.PieceBitboard(WhitePiece(Rook)).Count() - pos.PieceBitboard(BlackPiece(Rook)).Count())
	materialScore += QueenValue * (pos.PieceBitboard(WhitePiece(Queen)).Count() - pos.PieceBitboard(BlackPiece(Queen)).Count())

	return materialScore
}

func mobilityScore(pos *Position) int {
	c2m := pos.ColorToMove

	pos.ColorToMove = White
	whiteMobility := LegalMoves(pos).Count

	pos.ColorToMove = Black
	blackMobility := LegalMoves(pos).Count

	pos.ColorToMove = c2m

	return MobilityValue * (whiteMobility - blackMobility)
}

func who2move(c2m Color) int {
	if c2m == White {
		return 1
	}
	return -1
}

func PieceValue(piece PieceType) int {
	switch piece {
	case Pawn:
		return PawnValue
	case Knight:
		return KnightValue
	case Bishop:
		return BishopValue
	case Rook:
		return RookValue
	case Queen:
		return QueenValue
	default:
		return 0
	}
}
