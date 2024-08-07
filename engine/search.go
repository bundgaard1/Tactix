package engine

func Search(pos *Position, depth int) Move {

	moves := LegalMoves(pos)

	bestMove := moves.Moves[0]

	return bestMove
}
