package main

import (
	"fmt"
	"tactix/engine"
)

func main() {
	fmt.Println("Tactix!!")

	var pos engine.Position

	pos.FromStandardStartingPosition()

	pos.PrintPosition()

	n := Perft(&pos, 2)

	fmt.Printf("Perft nodes: %d\n", n)
}

func Perft(pos *engine.Position, depth int) int {
	nodes := 0

	moveList := pos.GenerateMoves()

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
