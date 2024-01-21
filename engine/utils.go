package engine

func Rank(square Square) int8 {
	return int8(square-1)/8 + 1
}

var FileRune = map[int8]rune{
	1: 'a',
	2: 'b',
	3: 'c',
	4: 'd',
	5: 'e',
	6: 'f',
	7: 'g',
	8: 'h',
}

func File(square Square) int8 {
	return int8(square-1)%8 + 1
}

func CastelingRightsToString(rights uint8) string {
	castingString := ""
	if rights == 0 {
		castingString = "-"
	} else {
		if rights&WhiteKingsideRight != 0 {
			castingString += "K"
		}
		if rights&WhiteQueensideRight != 0 {
			castingString += "Q"
		}
		if rights&BlackKingsideRight != 0 {
			castingString += "k"
		}
		if rights&BlackQueensideRight != 0 {
			castingString += "q"
		}
	}
	return castingString
}

func AlgToMove(algMove string) Move {
	move := Move{}
	from := algMove[:2]
	to := algMove[2:]

	fromFile := from[0] - 'a' + 1
	move.From = Square((from[1]-'0'-1)*8 + fromFile)

	toFile := to[0] - 'a' + 1
	move.To = Square((from[1]-'0'-1)*8 + toFile)

	return move

}
