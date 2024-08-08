package engine_test

import (
	"fmt"
	"tactix/engine"
	"testing"
)

var bitboardsForStandard map[engine.Piece]engine.Bitboard = map[engine.Piece]engine.Bitboard{
	engine.WhitePiece(engine.Pawn):   0x000000000000FF00, // White Pawns
	engine.WhitePiece(engine.Knight): 0x0000000000000042, // White Knights
	engine.WhitePiece(engine.Bishop): 0x0000000000000024, // White Bishops
	engine.WhitePiece(engine.Rook):   0x0000000000000081, // White Rooks
	engine.WhitePiece(engine.Queen):  0x0000000000000008, // White Queens
	engine.WhitePiece(engine.King):   0x0000000000000010, // White Kings

	engine.BlackPiece(engine.Pawn):   0x00FF000000000000, // Black Pawns
	engine.BlackPiece(engine.Knight): 0x4200000000000000, // Black Knights
	engine.BlackPiece(engine.Bishop): 0x2400000000000000, // Black Bishops
	engine.BlackPiece(engine.Rook):   0x8100000000000000, // Black Rooks
	engine.BlackPiece(engine.Queen):  0x0800000000000000, // Black Queens
	engine.BlackPiece(engine.King):   0x1000000000000000, // Black Kings
}

// Test bitboard string representation
func TestPieceBitboardsFromStart(t *testing.T) {
	pos := engine.FromStandardStartingPosition()
	pos.InitPieceBitboards()

	for piece, bitboard := range bitboardsForStandard {
		if *pos.PieceBitboard(piece) != bitboard {
			t.Errorf("Bitboard for %s is incorrect", piece)
		}
	}
}

func TestPieceBitboardsMakeMove(t *testing.T) {
	for _, perftTest := range engine.PerftSuite {

		pos := engine.FromFEN(perftTest.FEN)

		moves := engine.LegalMoves(&pos)

		for i := 0; i < moves.Count; i++ {
			move := moves.Moves[i]

			pos.MakeMove(move)
			result, piece := positionBitboardsCorrect(&pos)
			if !result {
				fmt.Print(pos.String())
				fmt.Print(piece.String(), " : ", fmt.Sprintf(" %b \n", *pos.PieceBitboard(piece)))
				t.Error("Bitboards incorrect")
			}

			pos.UndoMove(move)

		}
	}
}

func positionBitboardsCorrect(pos *engine.Position) (bool, engine.Piece) {
	for i := engine.Square(1); i <= 64; i++ {
		piece := pos.Board[i]
		if piece == engine.ANoPiece() {
			continue
		}
		if !pos.PieceBitboard(piece).IsSet(i) {
			return false, piece
		}
	}
	return true, engine.ANoPiece()
}
