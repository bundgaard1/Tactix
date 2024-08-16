package engine

import (
	"fmt"
)

const (
	SearchDepth = 6
)

type Search struct {
	pos           Position
	searchOver    bool
	nodesSearched int
	timer         Timer

	openingBook OpeningBook
	opening     bool
}

func NewSearch(pos *Position) (search Search) {
	return Search{
		pos:           *pos,
		searchOver:    false,
		nodesSearched: 0,
		timer:         NewTimer(),
		openingBook:   *NewOpeningBook(),
		opening:       true,
	}
}

// Iterative deepening
func (search *Search) Search() Move {
	bestMove, bestScore := Move{}, NegativeInfinity

	for depth := 1; depth <= SearchDepth; depth++ {
		move, score := search.rootAlphaBeta(depth)

		if search.searchOver {
			return bestMove
		}

		if score > bestScore {
			bestMove, bestScore = move, score
		}

		search.searchInfo(depth, bestScore, bestMove)
	}

	return bestMove
}

func (search *Search) rootAlphaBeta(depth int) (Move, int) {
	alpha, beta := NegativeInfinity, PositiveInfinity
	search.nodesSearched = 0

	bestMove := Move{}

	moves := LegalMoves(&search.pos)
	search.orderMoves(&moves)

	for i := 0; i < len(moves); i++ {
		move := moves[i]

		search.pos.MakeMove(move)
		score := -search.alphaBeta(-beta, -alpha, depth-1)
		search.pos.UndoMove(move)

		// fmt.Println(" Move: ", move.UCIString(), " Score: ", score)
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

func (search *Search) alphaBeta(alpha, beta, depthLeft int) int {
	search.nodesSearched++

	if depthLeft == 0 {
		return search.quiesce(alpha, beta)
	}
	bestValue := NegativeInfinity

	moves := LegalMoves(&search.pos)
	search.orderMoves(&moves)

	for i := 0; i < len(moves); i++ {
		move := moves[i]

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
	search.nodesSearched++
	stand_pat := Evaluate(&search.pos)

	if stand_pat >= beta {
		return beta
	}
	if alpha < stand_pat {
		alpha = stand_pat
	}

	moves := LegalMoves(&search.pos)
	search.orderMoves(&moves)

	for i := 0; i < len(moves); i++ {
		move := moves[i]
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

func (search *Search) orderMoves(moves *MoveList) {
	var scores []int

	for i := 0; i < len(*moves); i++ {
		scores = append(scores, scoreMove((*moves)[i], &search.pos))
	}

	// Sort Moves based on scores
	// Better scores first
	for i := 0; i < len(*moves)-1; i++ {
		for j := i + 1; j < len(*moves); j++ {
			if scores[j] > scores[i] {
				// swap
				scores[j], scores[i] = scores[i], scores[j]
				(*moves)[j], (*moves)[i] = (*moves)[i], (*moves)[j]
			}
		}
	}
}

func scoreMove(move Move, pos *Position) int {
	scoreGuess := 0

	movePieceType := pos.Board[move.From].PType
	capturedPieceType := pos.Board[move.To].PType

	if capturedPieceType != NoPiece {
		scoreGuess += 10*PieceValue(capturedPieceType) - PieceValue(movePieceType)
	}

	if move.Flag != NoFlag {
		scoreGuess += 100
	}

	if move.Flag.IsPromotion() {
		scoreGuess += PieceValue(Queen)
	}

	return scoreGuess
}

func (search *Search) searchInfo(depth int, bestScore int, bestMove Move) {
	fmt.Printf(
		"info depth %d score %d nodes %d bestmove %s\n",
		depth, bestScore,
		search.nodesSearched,
		bestMove.UCIString(),
	)
}
