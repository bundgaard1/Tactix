package engine_test

import (
	"fmt"
	"tactix/engine"
	"testing"
)

func TestPerftSuite(t *testing.T) {
	for i, perftTest := range engine.PerftSuite {
		pos, err := engine.FromFEN(perftTest.FEN)
		if err != nil {
			t.Error(err)
		}

		nodesExplored := engine.Perft(pos, perftTest.Depth)

		if nodesExplored != perftTest.ExpectedNodes {
			fmt.Println(engine.FEN(pos))
			t.Error("wrong at ", i, " : ", nodesExplored, " (expected: ", perftTest.ExpectedNodes, ") \n")
		}
	}
}
