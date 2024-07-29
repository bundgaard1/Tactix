package engine

import (
	"fmt"
	"time"
)

func Perft(pos *Position, depth int) int {
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

func PerftDivided(pos *Position, depth int) int {
	start := time.Now()

	nodes := 0

	moveList := GetValidMoves(pos)

	if depth == 1 {
		return moveList.Count
	}

	for i := 0; i < moveList.Count; i++ {
		move := moveList.Moves[i]
		pos.MakeMove(move)
		newNodes := Perft(pos, depth-1)
		pos.UndoMove(move)

		fmt.Printf("%s : %d \n", move.String(), newNodes)
		nodes += newNodes
	}

	duration := time.Since(start)

	fmt.Printf("\n Perft nodes: %d\n", nodes)
	fmt.Printf("	Runtime: %d ms \n", duration.Milliseconds())
	nodesPerSecond := nodes / int(duration.Milliseconds())
	fmt.Printf("	nodes/s: %d k \n\n ", nodesPerSecond)

	return nodes
}
