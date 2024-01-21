package main

import (
	"fmt"
	"tactix/engine"
	"time"
)

func main() {
	fmt.Println("Tactix!!")

	var pos engine.Position

	pos.FromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")

	// pos.MakeMove(engine.Move{From: 7, To: 24, Flag: engine.NoFlag})
	// pos.MakeMove(engine.Move{From: 49, To: 41, Flag: engine.NoFlag})

	pos.PrintPosition()

	PerftDivided(&pos, 5)

	pos.PrintPosition()

}

func Perft(pos *engine.Position, depth int) int {
	nodes := 0

	moveList := pos.GenerateLegalMoves()

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

func PerftDivided(pos *engine.Position, depth int) int {
	start := time.Now()

	nodes := 0

	moveList := pos.GenerateLegalMoves()

	if depth == 1 {
		// for i := 0; i < moveList.Count; i++ {
		// 	moveList.Moves[i].Print()
		// }
		return moveList.Count
	}

	for i := 0; i < moveList.Count; i++ {
		move := moveList.Moves[i]
		move.Print()
		pos.MakeMove(move)
		newNodes := Perft(pos, depth-1)
		pos.UndoMove(move)
		fmt.Printf("%d \n", newNodes)
		nodes += newNodes
	}

	duration := time.Since(start)

	fmt.Printf("\n Perft nodes: %d\n", nodes)
	fmt.Printf("	Runtime: %d ms \n", duration.Milliseconds())
	nodesPerSecond := nodes / int(duration.Milliseconds())
	fmt.Printf("	nodes/s: %d k \n\n ", nodesPerSecond)

	return nodes
}
