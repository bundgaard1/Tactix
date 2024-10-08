package engine_test

import (
	"fmt"
	"tactix/engine"
	"testing"
)

func TestBitboards(t *testing.T) {
	var bb engine.Bitboard

	for i := 1; i <= 64; i++ {

		bb.Set(engine.Square(i))
		if !bb.IsSet(engine.Square(i)) {
			t.Error("bit should be set: ", i)
			fmt.Print(bb.String())
		}

		bb.Unset(engine.Square(i))

		if bb.IsSet(engine.Square(i)) {
			t.Error("bit should be cleared: ", i)
			fmt.Print(bb.String())
		}
	}

	count := bb.Count()
	if count != 0 {
		t.Errorf("expected 0 bit, got %d", count)
	}

	bb = 0b1111
	count = bb.Count()
	if count != 4 {
		t.Errorf("expected 4 bit, got %d", count)
	}
}
