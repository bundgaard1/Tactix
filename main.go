package main

import (
	"fmt"
	"tactix/cli"
	"tactix/engine"
)

func main() {

	cli.Run()

	pos := engine.FromFEN("2k5/K6r/8/2Pp4/8/8/8/8 w - d6 0 6")
	fmt.Print(pos.String())
	bb := engine.BBKingAttackedMask(&pos)
	moves := engine.GetValidMoves(&pos)

	for i := 0; i < moves.Count; i++ {
		move := moves.Moves[i]
		fmt.Println(move.String())
	}

	fmt.Print(bb.StringOnBoard())

}
