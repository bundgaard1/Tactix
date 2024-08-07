package main

import (
	"fmt"
	"tactix/engine"
)

func main() {
	//cli.Run()

	pos := engine.FromStandardStartingPosition()
	// piece := engine.Piece{Color: engine.Black, PieceType: engine.King}

	fmt.Println(pos.ColorBitboard(engine.White).StringOnBoard())
	fmt.Println(pos.ColorBitboard(engine.Black).StringOnBoard())
	fmt.Println(pos.String())

}
