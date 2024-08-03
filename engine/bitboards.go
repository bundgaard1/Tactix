package engine

import (
	"fmt"
	"math/bits"
	"strings"
)

// BitBoard representation
type Bitboard uint64

const (
	FullBB  Bitboard = 0xffff_ffff_ffff_ffff
	EmptyBB Bitboard = 0x0
)

func (bitboard *Bitboard) SetBit(sq Square) {
	*bitboard |= (1 << (sq - 1))
}

func (bitboard *Bitboard) ClearBit(sq Square) {
	*bitboard &= FullBB ^ (1 << (sq - 1))
}

func (bitboard Bitboard) IsBitSet(sq Square) bool {
	return (bitboard & (1 << (sq - 1))) != 0
}

// Least significant bit
func (bitboard Bitboard) Lsb() Square {
	return Square(bits.TrailingZeros64(uint64(bitboard)))
}

func (bitboard *Bitboard) PopBit() Square {
	pos := bitboard.Lsb()
	bitboard.ClearBit(pos + 1)
	return pos + 1
}

// Count the bits in a given bitboard using the SWAR-popcount
// algorithm for 64-bit integers.
func (bitboard Bitboard) CountBits() int {
	return bits.OnesCount64(uint64(bitboard))
}

func (bb Bitboard) String() string {
	return fmt.Sprintf("%064b\n", bb)
}

func (bb Bitboard) StringOnBoard() string {
	var str strings.Builder
	for row := 7; row >= 0; row-- {
		for file := 1; file <= 8; file++ {
			i := row*8 + file - 1
			value := (bb >> i) & 1
			if value == 1 {
				str.WriteString("x ")
			} else {
				str.WriteString(". ")
			}
		}
		str.WriteString("\n")
	}
	return str.String()
}
