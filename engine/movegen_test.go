package engine_test

import (
	"tactix/engine"
	"testing"
)

func TestStartingPosition(t *testing.T) {
	var pos engine.Position
	pos.FromStandardStartingPosition()

	var moveList engine.MoveList = pos.GenerateMoves()

	if moveList.Count != 20 {
		t.Errorf("Result was incorrect, got: %d, want: %d.", moveList.Count, 20)
	}

}
