package engine

import "fmt"

const (
	SearchDepth = 6
)

type Search struct {
	pos           Position
	searchOver    bool
	nodesSearched int
}

func NewSearch(pos *Position) (search Search) {
	search.pos = *pos
	search.searchOver = false

	return search
}

// Iterative deepening
func (search *Search) Search() Move {
	bestMove, bestScore := Move{}, NegativeInfinity

	for depth := 1; depth < SearchDepth; depth++ {
		move, score := search.rootAlphaBeta(SearchDepth)
		if search.searchOver {
			return bestMove
		}

		bestMove, bestScore = move, score

		search.searchInfo(depth, bestScore)
	}

	return bestMove
}

func (search *Search) rootAlphaBeta(depth int) (Move, int) {
	alpha, beta := NegativeInfinity, PositiveInfinity

	bestMove := Move{}

	moves := LegalMoves(&search.pos)

	for i := 0; i < moves.Count; i++ {
		move := moves.Moves[i]

		search.pos.MakeMove(move)
		score := -search.alphaBeta(-beta, -alpha, depth-1)
		search.pos.UndoMove(move)

		if score == PositiveInfinity {
			return move, beta
		}

		if score > alpha {
			alpha = score
			bestMove = move
		}

	}

	return bestMove, alpha
}

// alpha is the score of the max-player
// beta is the score of the min-player
func (search *Search) alphaBeta(alpha, beta, depthLeft int) int {
	if depthLeft == 0 {
		return search.quiesce(alpha, beta)
	}
	bestValue := NegativeInfinity

	moves := LegalMoves(&search.pos)
	for i := 0; i < moves.Count; i++ {
		move := moves.Moves[i]

		search.pos.MakeMove(move)
		score := -search.alphaBeta(-beta, -alpha, depthLeft-1)
		search.pos.UndoMove(move)

		if score > bestValue {
			bestValue = score
			if score > alpha {
				alpha = score
			}
		}
		if score >= beta {
			return bestValue
		}
	}
	return bestValue
}

func (search *Search) quiesce(alpha, beta int) int {
	stand_pat := Evaluate(&search.pos)

	if stand_pat >= beta {
		return beta
	}
	if alpha < stand_pat {
		alpha = stand_pat
	}

	moves := LegalMoves(&search.pos)

	for i := 0; i <= moves.Count; i++ {
		move := moves.Moves[i]
		if !search.pos.isCapture(move) {
			continue
		}
		search.pos.MakeMove(move)
		score := -search.quiesce(-beta, -alpha)
		search.pos.UndoMove(move)

		if score >= beta {
			return beta
		}
		if score > alpha {
			alpha = score
		}
	}

	return alpha
}

func (pos *Position) isCapture(move Move) bool {
	if move.Flag == EnPassentCapture {
		return true
	}
	return pos.Board[move.To].Color == pos.ColorToMove.opposite()
}

func (search *Search) searchInfo(depth, bestScore int) {
	fmt.Printf(
		"info depth %d score cp %d nodes %d\n",
		depth, bestScore,
		search.nodesSearched,
	)
}
