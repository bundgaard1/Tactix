package engine

type MoveFlag uint8

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
	From int8
	To   int8
	Flag MoveFlag
}
