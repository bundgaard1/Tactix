package engine

func (pos *Position) GenerateLegalMoves() MoveList {
	// ------- NAIVE SOLUTION --------
	// Very Slow
	var pseudoLegalMoves MoveList = pos.GenerateMoves()
	var legalMoves MoveList

	for i := 0; i < pseudoLegalMoves.Count; i++ {
		moveToVarify := pseudoLegalMoves.Moves[i]
		pos.MakeMove(moveToVarify)
		var oppoenentResponses = pos.GenerateMoves()
		var addMove bool = true

		for j := 0; j < oppoenentResponses.Count; j++ {
			oppoRespons := oppoenentResponses.Moves[j]
			targetPiece := pos.Board[oppoRespons.To]
			if targetPiece.Color == pos.ColorToMove.opposite() && targetPiece.PieceType == King {
				addMove = false
				break
			}
		}
		if addMove {
			legalMoves.AddMove(moveToVarify)
		}
		pos.UndoMove(moveToVarify)
	}

	return legalMoves
}

func (pos *Position) GenerateLegalMoves2() MoveList {
	// var pseudoLegalMoves MoveList = pos.GenerateMoves()
	var legalMoves MoveList

	// Kingmoves
	// 		Calculate attacked square,
	//		and if it is a king that moves, make sure it is not into an attacked square

	// Check evasions
	//  	Single checks
	// 		Double checks

	return legalMoves
}

// Generate Pseudo-Legal Moves
func (pos *Position) GenerateMoves() MoveList {

	var moves MoveList

	for i := Square(1); i <= 64; i++ {
		currPiece := pos.Board[i]
		if currPiece.Color != pos.ColorToMove {
			continue
		}

		switch currPiece.PieceType {
		case Pawn:
			moves.AddMoves(genPawnMoves(i, currPiece, pos))
		case Knight:
			moves.AddMoves(GenKnightMoves(i, currPiece, pos))
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

func genPawnMoves(square Square, piece Piece, pos *Position) MoveList {
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

func GenKnightMoves(square Square, piece Piece, pos *Position) MoveList {

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

	// fmt.Printf("%d : file: %d  Rank:  %d  =  %08b\n", square, File(square), Rank(square), allowedMoves)

	for i := 0; i < 8; i++ {
		// fmt.Printf("%d : %d\n", i, (allowedMoves>>i)&0b1)
		if ((allowedMoves>>i)&1) == 1 && (pos.Board[square+knightMoveOffsets[i]].Color != piece.Color) {
			knightMoveList.AddMove(Move{From: square, To: square + knightMoveOffsets[i], Flag: NoFlag})
		}

	}
	return knightMoveList
}

var directionOffsets = [8]Square{8, -8, 1, -1, 7, -7, 9, -9}

var SqToEdgeComuted bool = false
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

	SqToEdgeComuted = true
}

func GenSlidingMoves(square Square, piece Piece, pos *Position) MoveList {
	if !SqToEdgeComuted {
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

func GenKingMoves(square Square, piece Piece, pos *Position) MoveList {
	if !SqToEdgeComuted {
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

	// Castling
	if square == Square(5) && piece.Color == White {
		if pos.CastlingRights&WhiteKingsideRight != 0 && pos.Board[7].PieceType|pos.Board[6].PieceType == NoPiece {
			kingMoveList.AddMove(Move{From: square, To: square + 2, Flag: Castling})
		}
		if pos.CastlingRights&WhiteQueensideRight != 0 && pos.Board[2].PieceType|pos.Board[3].PieceType|pos.Board[4].PieceType == NoPiece {
			kingMoveList.AddMove(Move{From: square, To: square - 2, Flag: Castling})
		}
	}
	if square == Square(61) && piece.Color == Black {
		if pos.CastlingRights&BlackKingsideRight != 0 && pos.Board[62].PieceType|pos.Board[63].PieceType == NoPiece {
			kingMoveList.AddMove(Move{From: square, To: square + 2, Flag: Castling})
		}
		if pos.CastlingRights&BlackQueensideRight != 0 && pos.Board[58].PieceType|pos.Board[59].PieceType|pos.Board[60].PieceType == NoPiece {
			kingMoveList.AddMove(Move{From: square, To: square - 2, Flag: Castling})
		}
	}

	return kingMoveList
}
