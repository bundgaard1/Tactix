package engine

func (pos *Position) GenerateMoves() MoveList {

	var moves MoveList
	var i int8
	for i = 1; i < 65; i++ {
		currPiece := pos.Board[i]
		if currPiece.Color != pos.ColorToMove {
			continue
		}

		switch currPiece.PieceType {
		case Pawn:
			moves.AddMoves(genPawnMoves(i, currPiece.Color, pos))
		case Knight:
			moves.AddMoves(GenKnightMoves(i, currPiece.Color, pos))
		case Bishop, Queen, Rook:
			moves.AddMoves(GenSlidingMoves(i, currPiece, pos))
		case King:
			moves.AddMoves(GenKingMoves(i, currPiece, pos))
		default:
			break
		}

	}
	return moves
}

func genPawnMoves(square int8, color Color, pos *Position) MoveList {
	var pawnMoveList MoveList
	var offset int8
	if color == White {
		offset = 8
	} else {
		offset = -8
	}
	// Move forward
	if pos.Board[square+offset].PieceType == NoPiece {
		pawnMoveList.AddMove(Move{From: square, To: square + offset, Flag: NoFlag})
	}
	push := square + 2*offset
	// Push
	if ((color == White && Rank(square) == 2) ||
		(color == Black && Rank(square) == 7)) && pos.Board[square+offset].PieceType == NoPiece {
		pawnMoveList.AddMove(Move{From: square, To: push, Flag: NoFlag})
	}

	// Attack
	if File(square) != 1 && pos.Board[square+offset-1].PieceType != NoPiece && pos.Board[square+offset-1].Color != color {
		pawnMoveList.AddMove(Move{From: square, To: square + offset - 1, Flag: NoFlag})
	}
	if File(square) != 8 && pos.Board[square+offset+1].PieceType != NoPiece && pos.Board[square+offset+1].Color != color {
		pawnMoveList.AddMove(Move{From: square, To: square + offset + 1, Flag: NoFlag})
	}

	// Promotion
	if (Rank(square) == 7 && color == White) || (Rank(square) == 1 && color == Black) {
		var countBefore int = pawnMoveList.Count

		for i := 0; i < countBefore; i++ {
			pawnMoveList.Moves[i].Flag = PromotionToQueen
			var to int8 = pawnMoveList.Moves[i].To
			pawnMoveList.AddMove(Move{square, to, PromotionToBishop})
			pawnMoveList.AddMove(Move{square, to, PromotionToQueen})
			pawnMoveList.AddMove(Move{square, to, PromotionToKnight})
		}
	}

	// Still missing En passant

	return pawnMoveList
}

var knightMoveOffsets = [8]int8{-17, -15, -6, -10, 6, 10, 15, 17}

func GenKnightMoves(square int8, color Color, pos *Position) MoveList {

	var allowedMoves uint8 = 0b1111_1111
	switch Rank(square) {
	case 1:
		allowedMoves &= 0b1111_0000
	case 2:
		allowedMoves &= 0b1111_1100
	case 7:
		allowedMoves &= 0b0011_1111
	case 8:
		allowedMoves &= 0b0000_1111
	}
	switch File(square) {
	case 1:
		allowedMoves &= 0b1010_0101
	case 2:
		allowedMoves &= 0b1110_0111
	case 7:
		allowedMoves &= 0b1101_1011
	case 8:
		allowedMoves &= 0b0101_1010
	}

	var knightMoveList MoveList

	for i := 0; i < 8; i++ {
		if ((allowedMoves>>i)&0b1) == 0b1 && (pos.Board[square+knightMoveOffsets[i]].Color != color) {
			knightMoveList.AddMove(Move{From: square, To: square + knightMoveOffsets[i], Flag: NoFlag})
		}
	}
	return knightMoveList
}

var directionOffsets = [8]int8{8, -8, 1, -1, 7, -7, 9, -9}

var SqToEdgeComuted bool = false
var numSquaresToEdge [65][8]int8

func computeSquaresToEdge() {

	for i := 1; i < 65; i++ {
		numUp := 8 - Rank(int8(i))
		numDown := Rank(int8(i)) - 1
		numLeft := File(int8(i)) - 1
		numRight := 8 - File(int8(i))

		numSquaresToEdge[i][0] = numUp
		numSquaresToEdge[i][1] = numDown
		numSquaresToEdge[i][2] = numRight
		numSquaresToEdge[i][3] = numLeft
		numSquaresToEdge[i][4] = min(numUp, numLeft)
		numSquaresToEdge[i][5] = min(numDown, numRight)
		numSquaresToEdge[i][6] = min(numUp, numRight)
		numSquaresToEdge[i][7] = min(numDown, numLeft)
	}

	SqToEdgeComuted = true
}

func GenSlidingMoves(square int8, piece Piece, pos *Position) MoveList {
	if !SqToEdgeComuted {
		computeSquaresToEdge()
	}

	var slidingMoveList MoveList

	var dirStartIndex int8 = 0
	var dirEndIndex int8 = 8

	if piece.PieceType == Bishop {
		dirStartIndex = 4
	}
	if piece.PieceType == Rook {
		dirEndIndex = 4
	}

	for dirIndex := dirStartIndex; dirIndex < dirEndIndex; dirIndex++ {
		var i int8
		for i = 0; i < numSquaresToEdge[square][dirIndex]; i++ {
			var to int8 = square + directionOffsets[dirIndex]*(i+1)

			var toColor Color = pos.Board[to].Color
			if toColor == piece.Color {
				break
			}

			slidingMoveList.AddMove(Move{From: square, To: to, Flag: NoFlag})

			if toColor == piece.Color.opposite() {
				break
			}
		}
	}

	return slidingMoveList
}

func GenKingMoves(square int8, piece Piece, pos *Position) MoveList {
	if !SqToEdgeComuted {
		computeSquaresToEdge()
	}

	var kingMoveList MoveList

	for dirIndex := 0; dirIndex < 8; dirIndex++ {

		if numSquaresToEdge[square][dirIndex] < 1 {
			continue
		}

		var to int8 = square + directionOffsets[dirIndex]
		var pieceOnTo Piece = pos.Board[to]

		if pieceOnTo.Color == piece.Color {
			continue
		}
		kingMoveList.AddMove(Move{From: square, To: to, Flag: NoFlag})
	}

	return kingMoveList
}
