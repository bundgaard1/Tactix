package engine

// ------- NAIVE SOLUTION --------
// Very Slow
func GetValidMoves(pos *Position) MoveList {

	var legalMoves MoveList

	return legalMoves
}

// Generate Pseudo-Legal Moves
func GetAllPossibleMoves(pos *Position) MoveList {

	var moves MoveList

	for sq := Square(1); sq <= 64; sq++ {
		currPiece := pos.Board[sq]
		if currPiece.Color != pos.ColorToMove {
			continue
		}

		switch currPiece.PieceType {
		case Pawn:
			moves.AddMoves(genPawnMoves(pos, sq, currPiece))
		case Knight:
			moves.AddMoves(genKnightMoves(pos, sq, currPiece))
		case Bishop, Queen, Rook:
			moves.AddMoves(genSlidingMoves(pos, sq, currPiece))
		case King:
			moves.AddMoves(genKingMoves(pos, sq, currPiece))
		default:
			continue
		}

	}
	return moves
}

func genPawnMoves(pos *Position, square Square, piece Piece) MoveList {
	var pawnMoveList MoveList
	var offset Square
	if piece.Color == White {
		offset = 8
	} else {
		offset = -8
	}

	// Move forward
	if pos.Board[square+offset].PieceType == NoPiece {
		pawnMoveList.AddMove(Move{From: square, To: square + offset, Flag: NoFlag})
	}

	// Push
	push := square + 2*offset
	if ((piece.Color == White && Rank(square) == 2) ||
		(piece.Color == Black && Rank(square) == 7)) && pos.Board[square+offset].PieceType|pos.Board[push].PieceType == NoPiece {
		pawnMoveList.AddMove(Move{From: square, To: push, Flag: PawnPush})
	}

	// Attack
	if File(square) != 1 && pos.Board[square+offset-1].PieceType != NoPiece && pos.Board[square+offset-1].Color != piece.Color {
		pawnMoveList.AddMove(Move{From: square, To: square + offset - 1, Flag: NoFlag})
	}
	if File(square) != 8 && pos.Board[square+offset+1].PieceType != NoPiece && pos.Board[square+offset+1].Color != piece.Color {
		pawnMoveList.AddMove(Move{From: square, To: square + offset + 1, Flag: NoFlag})
	}

	// Promotion
	if (Rank(square) == 7 && piece.Color == White) || (Rank(square) == 2 && piece.Color == Black) {
		var countBefore int = pawnMoveList.Count

		for i := 0; i < countBefore; i++ {
			pawnMoveList.Moves[i].Flag = PromotionToQueen
			var to Square = pawnMoveList.Moves[i].To
			pawnMoveList.AddMove(Move{square, to, PromotionToBishop})
			pawnMoveList.AddMove(Move{square, to, PromotionToQueen})
			pawnMoveList.AddMove(Move{square, to, PromotionToKnight})
			pawnMoveList.AddMove(Move{square, to, PromotionToRook})
		}
	}

	if pos.EPFile != 0 && ((Rank(square) == 5 && piece.Color == White) || (Rank(square) == 4 && piece.Color == Black)) {
		if File(square)-1 == pos.EPFile {
			pawnMoveList.AddMove(Move{From: square, To: square + offset - 1, Flag: EnPassentCapture})
		}
		if File(square)+1 == pos.EPFile {
			pawnMoveList.AddMove(Move{From: square, To: square + offset + 1, Flag: EnPassentCapture})
		}
	}

	return pawnMoveList
}

var knightMoveOffsets = [8]Square{-17, -15, -10, -6, 6, 10, 15, 17}

func genKnightMoves(pos *Position, square Square, piece Piece) MoveList {

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
		allowedMoves &= 0b1010_1010
	case 2:
		allowedMoves &= 0b1110_1011
	case 7:
		allowedMoves &= 0b1101_0111
	case 8:
		allowedMoves &= 0b0101_0101
	}

	var knightMoveList MoveList

	for i := 0; i < 8; i++ {
		if ((allowedMoves>>i)&1) == 1 && (pos.Board[square+knightMoveOffsets[i]].Color != piece.Color) {
			knightMoveList.AddMove(Move{From: square, To: square + knightMoveOffsets[i], Flag: NoFlag})
		}

	}
	return knightMoveList
}

var directionOffsets = [8]Square{8, -8, 1, -1, 7, -7, 9, -9}

var SqToEdgeComputed bool = false
var numSquaresToEdge [65][8]int8

func computeSquaresToEdge() {
	for i := Square(1); i < 65; i++ {
		numUp := 8 - Rank(i)
		numDown := Rank(i) - 1
		numLeft := File(i) - 1
		numRight := 8 - File(i)

		numSquaresToEdge[i][0] = numUp
		numSquaresToEdge[i][1] = numDown
		numSquaresToEdge[i][2] = numRight
		numSquaresToEdge[i][3] = numLeft
		numSquaresToEdge[i][4] = min(numUp, numLeft)
		numSquaresToEdge[i][5] = min(numDown, numRight)
		numSquaresToEdge[i][6] = min(numUp, numRight)
		numSquaresToEdge[i][7] = min(numDown, numLeft)
	}

	SqToEdgeComputed = true
}

func genSlidingMoves(pos *Position, square Square, piece Piece) MoveList {
	if !SqToEdgeComputed {
		computeSquaresToEdge()
	}

	var slidingMoveList MoveList

	var dirStartIndex Square = 0
	var dirEndIndex Square = 8

	if piece.PieceType == Bishop {
		dirStartIndex = 4
	}
	if piece.PieceType == Rook {
		dirEndIndex = 4
	}

	for dirIndex := dirStartIndex; dirIndex < dirEndIndex; dirIndex++ {
		for i := int8(0); i < numSquaresToEdge[square][dirIndex]; i++ {
			var to Square = square + directionOffsets[dirIndex]*Square(i+1)

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

func genKingMoves(pos *Position, square Square, piece Piece) MoveList {
	if !SqToEdgeComputed {
		computeSquaresToEdge()
	}

	var kingMoveList MoveList

	// Moving
	for dirIndex := 0; dirIndex < 8; dirIndex++ {
		if numSquaresToEdge[square][dirIndex] < 1 {
			continue
		}

		to := square + directionOffsets[dirIndex]
		pieceOnTo := pos.Board[to]

		if pieceOnTo.Color == piece.Color {
			continue
		}
		kingMoveList.AddMove(Move{From: square, To: to, Flag: NoFlag})
	}

	return kingMoveList
}

func genCastlingMoves(pos *Position, square Square, piece Piece) MoveList {
	var castlingMoves MoveList

	if squareUnderAttack(pos, square) {
		return MoveList{}
	}

	wk, wq, bk, bq := pos.getCastlingRights()

	if square == Square(E1) && piece.Color == White {
		if wk && pos.Board[F1].PieceType|pos.Board[G1].PieceType == NoPiece &&
			!squareUnderAttack(pos, F1) && !squareUnderAttack(pos, G1) {
			castlingMoves.AddMove(Move{From: square, To: square + 2, Flag: Castling})
		}
		if wq && pos.Board[2].PieceType|pos.Board[3].PieceType|pos.Board[4].PieceType == NoPiece &&
			!squareUnderAttack(pos, D1) && !squareUnderAttack(pos, C1) {
			castlingMoves.AddMove(Move{From: square, To: square - 2, Flag: Castling})
		}
	}
	if square == Square(E8) && piece.Color == Black {
		if bk && pos.Board[F8].PieceType|pos.Board[G8].PieceType == NoPiece &&
			!squareUnderAttack(pos, F8) && !squareUnderAttack(pos, G8) {
			castlingMoves.AddMove(Move{From: square, To: square + 2, Flag: Castling})
		}
		if bq && pos.Board[B8].PieceType|pos.Board[C8].PieceType|pos.Board[D8].PieceType == NoPiece &&
			!squareUnderAttack(pos, D8) && !squareUnderAttack(pos, C8) {
			castlingMoves.AddMove(Move{From: square, To: square - 2, Flag: Castling})
		}
	}

	return castlingMoves
}

func squareUnderAttack(pos *Position, sq Square) bool {
	if !SqToEdgeComputed {
		computeSquaresToEdge()
	}

	pos.ColorToMove = pos.ColorToMove.opposite()
	opponentMoves := GetAllPossibleMoves(pos)
	pos.ColorToMove = pos.ColorToMove.opposite()

	for i := 0; i < opponentMoves.Count; i++ {
		if opponentMoves.Moves[i].To == sq {
			return true
		}
	}
	return false
}
