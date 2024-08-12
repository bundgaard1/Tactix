package engine_test

import (
	"tactix/engine"
	"testing"
)

func TestMovelist(t *testing.T) {

	ml := engine.NewMoveList()

	ml.Append(engine.Move{From: 0, To: 1, Flag: engine.NoFlag})

	if len(*ml) != 1 {
		t.Errorf("Expected 1 move, got %d", len(*ml))
	}

	// Add many
	for i := 0; i < 100; i++ {
		ml.Append(engine.Move{From: 0, To: 1, Flag: engine.NoFlag})
	}

	if len(*ml) != 101 {
		t.Errorf("Expected 101 moves, got %d", len(*ml))
	}
}
