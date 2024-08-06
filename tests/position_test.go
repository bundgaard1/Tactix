package engine_test

import (
	"tactix/engine"
	"testing"
)

func TestFromFEN1(t *testing.T) {
	pos := engine.FromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")

	if pos.ColorToMove != engine.White {
		t.Errorf("Color to move is incorrect")
	}

	if pos.CastlingRights != engine.WhiteKingsideRight|engine.WhiteQueensideRight|engine.BlackKingsideRight|engine.BlackQueensideRight {
		t.Errorf("Castling rights are incorrect")
	}

	if pos.EPFile != 0 {
		t.Errorf("En passant file is incorrect")
	}

	if pos.Rule50 != 0 {
		t.Errorf("Rule 50 is incorrect")
	}

	if pos.Ply != 1 {
		t.Errorf("Ply is incorrect")
	}

	if pos.Board[engine.A1].PieceType != engine.Rook || pos.Board[engine.A1].Color != engine.White {
		t.Errorf("A1 is incorrect")
	}

	if pos.Board[engine.B1].PieceType != engine.Knight || pos.Board[engine.B1].Color != engine.White {
		t.Errorf("B1 is incorrect")
	}
}

func TestFromFEN2(t *testing.T) {
	pos := engine.FromFEN("rnbqkbnr/pp3ppp/8/2pPp3/P1P3N1/1P6/3PKPPP/RNB2B1R b Kq b4 0 6")

	if pos.ColorToMove != engine.Black {
		t.Errorf("Color to move is incorrect")
	}

	if pos.CastlingRights != engine.WhiteKingsideRight|engine.BlackQueensideRight {
		t.Errorf("Castling rights are incorrect")
	}

	if pos.EPFile != 2 {
		t.Errorf("En passant file is incorrect")
	}

	if pos.Rule50 != 0 {
		t.Errorf("Rule 50 is incorrect")
	}

	if pos.Ply != 6 {
		t.Errorf("Ply is incorrect")
	}
}
