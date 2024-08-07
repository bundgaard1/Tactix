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

func BBFromSquares(squares ...Square) Bitboard {
	var bb Bitboard
	for _, sq := range squares {
		bb.Set(sq)
	}
	return bb
}

func (bitboard *Bitboard) Set(sq Square) {
	*bitboard |= (1 << (sq - 1))
}

func (bitboard *Bitboard) Unset(sq Square) {
	*bitboard &= FullBB ^ (1 << (sq - 1))
}

func (bitboard Bitboard) IsSet(sq Square) bool {
	return (bitboard & (1 << (sq - 1))) != 0
}

// Least significant bit
func (bitboard Bitboard) lsb() Square {
	return Square(bits.TrailingZeros64(uint64(bitboard)))
}

func (bitboard *Bitboard) Pop() Square {
	if *bitboard == 0 {
		panic("Pop called on empty bitboard,")
	}

	pos := bitboard.lsb()

	return pos + 1
}

// Count the bits in a given bitboard using the SWAR-popcount
// algorithm for 64-bit integers.
func (bitboard Bitboard) Count() int {
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
