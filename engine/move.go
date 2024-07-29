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

func (m Move) String() string {
	return fmt.Sprintf("From: %d, To: %d, Flag: %d", m.From, m.To, m.Flag)
}

func (m Move) UCIString() string {
	return fmt.Sprintf("%c%d%c%d", FileRune[File(m.From)], Rank(m.From), FileRune[File(m.To)], Rank(m.To))
}
