package engine_test

import (
	"fmt"
	"tactix/engine"
	"testing"
	"time"
)

func TestPerftSuite(t *testing.T) {
	for i, perftTest := range engine.PerftSuite {
		pos := engine.FromFEN(perftTest.FEN)

		nodesExplored := engine.Perft(&pos, perftTest.Depth)

		if nodesExplored != perftTest.ExpectedNodes {
			t.Error("wrong at ", i, " : ", nodesExplored, " (expected: ", perftTest.ExpectedNodes, ") \n")
		}
	}
}

func PerftWithBenchmark() {
	totalNodes := 0

	startTime := time.Now()

	for i, perftTest := range engine.PerftSuite {
		pos := engine.FromFEN(perftTest.FEN)

		nodesExplored := engine.Perft(&pos, perftTest.Depth)
		totalNodes += nodesExplored

		fmt.Print("Test ", i, ": ", nodesExplored, " ")
		if nodesExplored == perftTest.ExpectedNodes {
			fmt.Print("check \n")
		} else {
			fmt.Print("wrong (expected ", perftTest.ExpectedNodes, ")\n")
		}
	}
	duration := time.Since(startTime)
	fmt.Printf(("\n"))
	fmt.Printf("Nodes : %d \n", totalNodes)
	fmt.Printf("Time  : %v \n", duration)
	fmt.Printf("Speed : %.2f MN/s \n", float64(totalNodes)/duration.Seconds()/1_000_000)
}
