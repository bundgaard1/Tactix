package engine_test

import (
	"fmt"
	"tactix/engine"
	"testing"
)

func TestPerftSuite(t *testing.T) {

	for i, perftTest := range engine.PerftSuite {
		pos := engine.FromFEN(perftTest.FEN)

		nodesExplored := engine.Perft(&pos, perftTest.Depth)

		fmt.Print("Test ", i, ": ", nodesExplored, " ")
		if nodesExplored == perftTest.ExpectedNodes {
			fmt.Print("check \n")
		} else {
			fmt.Print("wrong (expected ", perftTest.ExpectedNodes, ")\n")

		}
	}
}
