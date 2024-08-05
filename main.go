package main

import (
	"fmt"
	"tactix/engine"
)

func main() {
	// cli.Run()

	Position := engine.FromFEN("r3k3/1p3p2/p2q2p1/bn3P2/1N2PQP1/PB6/3K1R1r/3R4 w - - 0 1")

	fmt.Print(Position.String())

	bb := engine.BBKingAttackedMask(&Position)

	fmt.Print(bb.StringOnBoard())
}
