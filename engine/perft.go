package engine

import (
	"fmt"
	"strings"
)

func Perft(pos *Position, depth int) int {
	if depth == 0 {
		return 1
	}
	nodes := 0

	moveList := GetValidMoves(pos)

	if depth == 1 {
		return moveList.Count
	}

	for i := 0; i < moveList.Count; i++ {
		pos.MakeMove(moveList.Moves[i])
		nodes += Perft(pos, depth-1)
		pos.UndoMove(moveList.Moves[i])
	}

	return nodes
}

// Returns a summary with the move and the number of nodes for each move
func PerftDivided(pos *Position, depth int) (string, int) {
	var str strings.Builder

	totalNodes := 0
	moveList := GetValidMoves(pos)

	for i := 0; i < moveList.Count; i++ {
		move := moveList.Moves[i]
		pos.MakeMove(move)
		newNodes := Perft(pos, depth-1)
		pos.UndoMove(move)

		str.WriteString(fmt.Sprintf("%s : %d \n", move.UCIString(), newNodes))
		totalNodes += newNodes
	}

	return str.String(), totalNodes
}
