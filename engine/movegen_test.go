package engine_test

import (
	"fmt"
	"tactix/engine"
	"testing"
	"time"
)

func TestStartingPosition(t *testing.T) {
	var pos engine.Position

	pos.FromStandardStartingPosition()

	var moveList engine.MoveList = pos.GenerateMoves()

	if moveList.Count != 20 {
		t.Errorf("Result was incorrect, got: %d, want: %d.", moveList.Count, 20)
	}
}

func TestPerftFromStart(t *testing.T) {
	var pos engine.Position

	pos.FromStandardStartingPosition()

	perftDepth4Nodes := 197281

	nodesDepth4Result := perft(&pos, 4)

	if perftDepth4Nodes != nodesDepth4Result {
		t.Error("Perft result not as expected: ", perftDepth4Nodes, " got: ", nodesDepth4Result)
	}
}

type PerftTestData struct {
	Pos           string
	Depth         int
	ExpectedNodes int
}

// Suite from https://www.chessprogramming.net/perfect-perft/
var PerftSuite = []PerftTestData{
	// {"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", 6, 119060324},
	// {"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq -", 5, 193690690},
	// {"8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - -", 7, 178633661},
	// {"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", 6, 706045033},
	{"1k6/1b6/8/8/7R/8/8/4K2R b K - 0 1", 5, 1063513},
	{"3k4/3p4/8/K1P4r/8/8/8/8 b - - 0 1", 6, 1134888},
	{"8/8/4k3/8/2p5/8/B2P2K1/8 w - - 0 1", 6, 1015133},
	{"8/8/1k6/2b5/2pP4/8/5K2/8 b - d3 0 1", 6, 1440467},
	{"5k2/8/8/8/8/8/8/4K2R w K - 0 1", 6, 661072},
	{"3k4/8/8/8/8/8/8/R3K3 w Q - 0 1", 6, 803711},
	{"r3k2r/1b4bq/8/8/8/8/7B/R3K2R w KQkq - 0 1", 4, 1274206},
	{"r3k2r/8/3Q4/8/8/5q2/8/R3K2R b KQkq - 0 1", 4, 1720476},
	{"2K2r2/4P3/8/8/8/8/8/3k4 w - - 0 1", 6, 3821001},
	{"8/8/1P2K3/8/2n5/1q6/8/5k2 b - - 0 1", 5, 1004658},
	{"4k3/1P6/8/8/8/8/K7/8 w - - 0 1", 6, 217342},
	{"8/P1k5/K7/8/8/8/8/8 w - - 0 1", 6, 92683},
	{"K1k5/8/P7/8/8/8/8/8 w - - 0 1", 6, 2217},
	{"8/k1P5/8/1K6/8/8/8/8 w - - 0 1", 7, 567584},
	{"8/8/2k5/5q2/5n2/8/5K2/8 b - - 0 1", 4, 23527},
}

func TestPerftSuite(t *testing.T) {

	for i, perftTest := range PerftSuite {
		var pos engine.Position
		pos.FromFEN(perftTest.Pos)

		nodesExplored := perft(&pos, perftTest.Depth)

		fmt.Print("Test ", i, " : ", nodesExplored, " ")
		if nodesExplored == perftTest.ExpectedNodes {
			fmt.Print("check \n")
		} else {
			fmt.Print("wrong (expected ", perftTest.ExpectedNodes, ")\n")

		}
	}
}

func perft(pos *engine.Position, depth int) int {
	nodes := 0

	moveList := pos.GenerateLegalMoves()

	if depth == 1 {
		return moveList.Count
	}

	for i := 0; i < moveList.Count; i++ {
		pos.MakeMove(moveList.Moves[i])
		nodes += perft(pos, depth-1)
		pos.UndoMove(moveList.Moves[i])
	}

	return nodes
}

func perftDivided(pos *engine.Position, depth int) int {
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
		newNodes := perft(pos, depth-1)
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
