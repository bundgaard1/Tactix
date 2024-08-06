package main

import (
	"fmt"
	"tactix/engine"
)

func main() {

	// cli.Run()

	pos := engine.FromFEN("8/4k3/8/5q2/5n2/5K2/8/8 w - - 0 4")
	fmt.Print(pos.String())
	bb := engine.BBSquaresUnderAttack(&pos)
	moves := engine.GetValidMoves(&pos)

	for i := 0; i < moves.Count; i++ {
		move := moves.Moves[i]
		fmt.Println(move.String())
	}

	fmt.Print(bb.StringOnBoard())

}
