package engine_test

import (
	"tactix/engine"
	"testing"
)

func TestFENBuilder1(t *testing.T) {
	pos := engine.FromFEN(engine.StartingPositionFEN)
	fen := engine.FEN(&pos)
	pos2 := engine.FromFEN(fen)
	fen2 := engine.FEN(&pos2)

	if fen != fen2 {
		t.Errorf("FENs do not match: %s != %s", fen, fen2)
	}

}

func TestFENBuilder2(t *testing.T) {
	// Test a position with a random configurations

	fenToTest := "8/8/1P2K3/8/2n5/1q6/8/5k2 b - - 0 1"

	pos := engine.FromFEN(fenToTest)
	fen := engine.FEN(&pos)
	pos2 := engine.FromFEN(fen)
	fen2 := engine.FEN(&pos2)

	if fen != fen2 {
		t.Errorf("FENs do not match: %s != %s", fen, fen2)
	}
}
