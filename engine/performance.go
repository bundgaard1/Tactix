package engine

import (
	"fmt"
	"time"
)

func SearchPerformance1() {
	pos, err := FromFEN("r3k2r/p1ppqpb1/Bn2pnp1/3PN3/1p2P3/2N2Q2/PPPB1PpP/R3K2R w KQkq - 0 1")
	if err != nil {
		panic(err)
	}

	s := NewSearch()
	s.SetPosition(pos)
	s.Timer.NoTimeControl()

	start := time.Now()
	s.SearchDepth(4)

	elapsed := time.Since(start)

	fmt.Printf("Nodes searched: %d \nElapsed Time: %.5f s", s.NodesSearched(), elapsed.Seconds())

}
