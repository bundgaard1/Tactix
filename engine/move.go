package engine

type MoveFlag uint8

const (
	None MoveFlag = iota
	PawnPush
	Castling
	EnPassentCapture
	PromotionToQueen
	PromotionToKnight
	PromotionToRook
	PromotionToBishop
)

type Move struct {
	From uint8
	To   uint8
	Flag MoveFlag
}
