package bbcalc

import (
	"fmt"
	"tactix/engine"
)

var knightMoveOffsets = [8]engine.Square{-17, -15, -10, -6, 6, 10, 15, 17}

func CalculateKnightMoveBB() {

	for i := 1; i <= 64; i++ {
		var allowedMoves uint8 = 0b1111_1111
		switch engine.Rank(engine.Square(i)) {
		case 1:
			allowedMoves &= 0b1111_0000
		case 2:
			allowedMoves &= 0b1111_1100
		case 7:
			allowedMoves &= 0b0011_1111
		case 8:
			allowedMoves &= 0b0000_1111
		}
		switch engine.File(engine.Square(i)) {
		case 1:
			allowedMoves &= 0b1010_1010
		case 2:
			allowedMoves &= 0b1110_1011
		case 7:
			allowedMoves &= 0b1101_0111
		case 8:
			allowedMoves &= 0b0101_0101
		}

		var moveMask engine.BB = engine.EmptyBB

		fmt.Print(i, ": ")
		for j := 0; j < 8; j++ {
			if ((allowedMoves >> j) & 1) != 0 {
				destSquare := engine.Square(i) + knightMoveOffsets[j]
				// fmt.Print(destSquare, " ")
				moveMask.SetBit(destSquare)
			}
		}
		fmt.Printf("0x%x", moveMask)
		for moveMask.CountBits() > 0 {
			pos := moveMask.PopBit()
			fmt.Print(pos, " ")
		}

		fmt.Printf(",  \n")
	}
}
