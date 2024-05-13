package engine

import (
	"fmt"
	"math/bits"
)

// BitBoard representation
type BB uint64

const (
	FullBB  BB = 0xffff_ffff_ffff_ffff
	EmptyBB BB = 0x0
)

func (bitboard *BB) SetBit(sq Square) {
	*bitboard |= (1 << (sq - 1))
}

func (bitboard *BB) ClearBit(sq Square) {
	*bitboard &= FullBB ^ (1 << (sq - 1))
}

func (bitboard BB) IsBitSet(sq Square) bool {
	return (bitboard & (1 << (sq - 1))) != 0
}

// Least significant bit
func (bitboard BB) Lsb() Square {
	return Square(bits.TrailingZeros64(uint64(bitboard)))
}

func (bitboard *BB) PopBit() Square {
	pos := bitboard.Lsb()
	bitboard.ClearBit(pos + 1)
	return pos + 1
}

// Count the bits in a given bitboard using the SWAR-popcount
// algorithm for 64-bit integers.
func (bitboard BB) CountBits() int {
	return bits.OnesCount64(uint64(bitboard))
}

func (bb BB) Print() {
	fmt.Printf("%064b \n", bb)
}

func (bb BB) PrintOnBoard() {
	fmt.Printf("\n")
	for row := 7; row >= 0; row-- {
		for file := 1; file <= 8; file++ {
			i := row*8 + file - 1
			value := (bb >> (i)) & 1
			if value == 1 {
				fmt.Printf("x ")
			} else {
				fmt.Printf(". ")
			}
		}
		fmt.Printf("\n")
	}
}
