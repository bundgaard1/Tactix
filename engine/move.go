package engine

type (
	MoveFlag int8
	Square   int8
)

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

func (flag MoveFlag) IsPromotion() bool {
	return flag == PromotionToQueen || flag == PromotionToKnight || flag == PromotionToRook || flag == PromotionToBishop
}
