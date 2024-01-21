package engine

import "fmt"

type MoveFlag int8
type Square int8

const (
	NoFlag MoveFlag = iota
	PawnPush
	Castling
	EnPassentCapture
	PromotionToQueen
	PromotionToKnight
	PromotionToRook
	PromotionToBishop
)

type Move struct {
	From Square
	To   Square
	Flag MoveFlag
}

func (m Move) Print() {
	fmt.Printf("%c", FileRune[File(m.From)])
	fmt.Printf("%d", Rank(m.From))
	fmt.Printf("%c", FileRune[File(m.To)])
	fmt.Printf("%d ", Rank(m.To))

}
