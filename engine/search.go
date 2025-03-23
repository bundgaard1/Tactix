package engine

import (
	"fmt"
)

const (
	SearchDepth = 8
)

type Search struct {
	pos      *Position
	BestMove Move

	nodesSearched int

	Timer Timer

	Debug bool
}

func NewSearch() *Search {
	return &Search{
		pos:           NewPosition(),
		BestMove:      NilMove(),
		nodesSearched: 0,
		Timer:         NewTimer(),
		Debug:         true,
	}
}

func (s *Search) Stop() {
	s.Timer.Stop = true
	Log("Search stopping")
}

func (s *Search) Reset() {
	s.Timer.Stop = false
}

func (s *Search) SetPosition(pos *Position) {
	s.pos = pos
}

func (s *Search) Search() {
	bestMove, bestScore := Move{}, NegativeInfinity

	s.Timer.Start()
	Log(fmt.Sprintf("Max Search Time: %v", s.Timer.TimeForMove))

	for depth := 1; depth <= SearchDepth; depth++ {

		move, score := s.rootAlphaBeta(depth)

		s.Timer.Check()
		if s.Timer.Stop {
			s.BestMove = bestMove
			Log("Search stopped")
			return
		}

		if score > bestScore {
			bestMove, bestScore = move, score
		}

		if s.Debug {
			s.searchInfo(depth, bestScore, bestMove)
		}
	}

	s.BestMove = bestMove
}

func (s *Search) rootAlphaBeta(depth int) (Move, int) {
	alpha, beta := NegativeInfinity, PositiveInfinity
	s.nodesSearched = 0

	bestMove := Move{}

	moves := LegalMoves(s.pos)
	s.orderMoves(&moves)

	for i := 0; i < len(moves); i++ {
		move := moves[i]

		s.pos.MakeMove(move)
		score := -s.alphaBeta(-beta, -alpha, depth-1)
		s.pos.UndoMove(move)

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

func (s *Search) alphaBeta(alpha, beta, depthLeft int) int {
	s.nodesSearched++

	// Premature stop, if the time is up
	s.Timer.Check()
	if s.Timer.Stop {
		return NegativeInfinity
	}

	if depthLeft == 0 {
		return s.quiesce(alpha, beta)
	}
	bestValue := NegativeInfinity

	moves := LegalMoves(s.pos)
	s.orderMoves(&moves)

	for i := 0; i < len(moves); i++ {
		move := moves[i]

		s.pos.MakeMove(move)
		score := -s.alphaBeta(-beta, -alpha, depthLeft-1)
		s.pos.UndoMove(move)

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

func (s *Search) quiesce(alpha, beta int) int {
	s.nodesSearched++
	stand_pat := Evaluate(s.pos)

	if stand_pat >= beta {
		return beta
	}
	if alpha < stand_pat {
		alpha = stand_pat
	}

	moves := LegalMoves(s.pos)
	s.orderMoves(&moves)

	for i := 0; i < len(moves); i++ {
		move := moves[i]
		if !s.pos.isCapture(move) {
			continue
		}
		s.pos.MakeMove(move)
		score := -s.quiesce(-beta, -alpha)
		s.pos.UndoMove(move)

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
		scores = append(scores, scoreMove((*moves)[i], search.pos))
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
	message := fmt.Sprintf(
		"info depth %d score cp %d nodes %d bestmove %s",
		depth, bestScore,
		search.nodesSearched,
		bestMove.UCIString(),
	)
	fmt.Println(message)
	Log(message)
}
