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

func NilMove() Move {
	return Move{From: -1, To: -1, Flag: NoFlag}
}
