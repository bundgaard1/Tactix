package main

import (
	"fmt"
	"tactix/engine"
)

func main() {

	// cli.Run()

	pos := engine.FromFEN("3k4/8/8/K1Pp3r/8/8/8/8 w - d6 0 2")
	fmt.Print(pos.String())
	bb := engine.BBPinnedSquares(&pos)
	moves := engine.GetValidMoves(&pos)

	for i := 0; i < moves.Count; i++ {
		move := moves.Moves[i]
		fmt.Println(move.String())
	}

	fmt.Print(bb.StringOnBoard())

}
