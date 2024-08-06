package engine

func GetValidMoves(pos *Position) MoveList {
	var legalMoves MoveList
	checks := numChecks(pos)
	if checks > 0 {
		if checks == 1 { // One check, Move the king to a safe square, or block the attack
			allMoves := filterMovesToLegal(pos, GetAllPossibleMoves(pos))

			kingSquare := pos.GetKingSquare(pos.ColorToMove)
			KingAttackedLine := BBKingAttackedMask(pos)

			for i := 0; i < allMoves.Count; i++ {
				move := allMoves.Moves[i]
				if move.From == kingSquare {
					// King moves should be filtered from filterMovesToLegal
					legalMoves.AddMove(move)
				} else {
					// Block the attack/ Capture the attacking piece
					if KingAttackedLine.IsBitSet(move.To) {
						legalMoves.AddMove(move)
					}
					// En passent check, very disgusting.
					if move.Flag == EnPassentCapture && pos.EPFile == File(move.To) &&
						((KingAttackedLine.IsBitSet(move.To-8) && pos.ColorToMove == White) ||
							(KingAttackedLine.IsBitSet(move.To+8) && pos.ColorToMove == Black)) {
						legalMoves.AddMove(move)
					}
				}
			}
		} else { // Two or more checks,Move the king to a safe square
			kingMoves := genKingMoves(pos, pos.GetKingSquare(pos.ColorToMove), pos.Board[pos.GetKingSquare(pos.ColorToMove)])
			legalKingMoves := filterMovesToLegal(pos, kingMoves)
			legalMoves.AddMoves(legalKingMoves)
		}

	} else {
		allMoves := GetAllPossibleMoves(pos)
		legalMoves.AddMoves(filterMovesToLegal(pos, allMoves))
		legalMoves.AddMoves(genCastlingMoves(pos, pos.GetKingSquare(pos.ColorToMove), pos.Board[pos.GetKingSquare(pos.ColorToMove)]))
	}

	if legalMoves.Count == 0 {
		if checks > 0 {
			pos.checkmate = true
		} else {
			pos.stalemate = true
		}
	} else {
		pos.checkmate = false
		pos.stalemate = false
	}

	return legalMoves
}

func filterMovesToLegal(pos *Position, moves MoveList) MoveList {

	var legalMoves MoveList

	for i := 0; i < moves.Count; i++ {
		if isMoveLegal(pos, moves.Moves[i]) {
			legalMoves.AddMove(moves.Moves[i])
		}
	}

	return legalMoves

}

// Check for if move is involved in pin, or a king moves onto attacked square.
func isMoveLegal(pos *Position, move Move) bool {
	piece := pos.Board[move.From]

	if piece.Color != pos.ColorToMove {
		return false
	}
	// Pinned Piece
	if piece.PieceType != King && BBPinnedSquares(pos).IsBitSet(move.From) && !BBPinnedSquares(pos).IsBitSet(move.To) {
		return false
	}
	// King cant move to attacked square
	if piece.PieceType == King && BBSquaresUnderAttack(pos).IsBitSet(move.To) {
		return false
	}

	return true
}

// All Posible Moves, does check for pins.
// No castling.
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
	var pawnDirection Square
	if piece.Color == White {
		pawnDirection = 8
	} else {
		pawnDirection = -8
	}

	// Move forward
	if pos.Board[square+pawnDirection].PieceType == NoPiece {
		pawnMoveList.AddMove(Move{From: square, To: square + pawnDirection, Flag: NoFlag})
	}

	// Push
	push := square + 2*pawnDirection
	if ((piece.Color == White && Rank(square) == 2) ||
		(piece.Color == Black && Rank(square) == 7)) && pos.Board[square+pawnDirection].PieceType|pos.Board[push].PieceType == NoPiece {
		pawnMoveList.AddMove(Move{From: square, To: push, Flag: PawnPush})
	}

	// Attack
	attackMoves := genPawnAttackMoves(pos, square, piece)
	for i := 0; i < attackMoves.Count; i++ {
		if pos.Board[attackMoves.Moves[i].To].Color == piece.Color.opposite() {
			pawnMoveList.AddMove(attackMoves.Moves[i])
		}
	}

	// Promotion - If we are at the end of the board, add all possible promotions
	if (Rank(square) == 7 && piece.Color == White) || (Rank(square) == 2 && piece.Color == Black) {
		var countBefore int = pawnMoveList.Count
		for i := 0; i < countBefore; i++ {
			to := pawnMoveList.Moves[i].To
			pawnMoveList.Moves[i].Flag = PromotionToQueen
			pawnMoveList.AddMove(Move{square, to, PromotionToBishop})
			pawnMoveList.AddMove(Move{square, to, PromotionToKnight})
			pawnMoveList.AddMove(Move{square, to, PromotionToRook})
		}
	}

	// En passent
	if pos.EPFile != 0 && ((Rank(square) == 5 && piece.Color == White) || (Rank(square) == 4 && piece.Color == Black)) {
		if File(square)-1 == pos.EPFile {
			pawnMoveList.AddMove(Move{From: square, To: square + pawnDirection - 1, Flag: EnPassentCapture})
		}
		if File(square)+1 == pos.EPFile {
			pawnMoveList.AddMove(Move{From: square, To: square + pawnDirection + 1, Flag: EnPassentCapture})
		}
	}

	return pawnMoveList
}

// Does not check for promotion or if the attack has a target
func genPawnAttackMoves(pos *Position, square Square, piece Piece) MoveList {
	var pawnMoveList MoveList
	var pawnDirection Square
	if piece.Color == White {
		pawnDirection = 8
	} else {
		pawnDirection = -8
	}

	// Attack
	if File(square) != 1 {
		pawnMoveList.AddMove(Move{From: square, To: square + pawnDirection - 1, Flag: NoFlag})
	}
	if File(square) != 8 {
		pawnMoveList.AddMove(Move{From: square, To: square + pawnDirection + 1, Flag: NoFlag})
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
		if ((allowedMoves>>i)&1) == 1 && (IgnoreFriendlyPieces || pos.Board[square+knightMoveOffsets[i]].Color != piece.Color) {
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

	dirStartIndex, dirEndIndex := slidingStartAndEndIndex(piece.PieceType)

	for dirIndex := dirStartIndex; dirIndex < dirEndIndex; dirIndex++ {
		for i := int8(1); i <= numSquaresToEdge[square][dirIndex]; i++ {
			to := square + directionOffsets[dirIndex]*Square(i)

			toColor := pos.Board[to].Color
			if toColor == piece.Color {
				if IgnoreFriendlyPieces {
					slidingMoveList.AddMove(Move{From: square, To: to, Flag: NoFlag})
				}
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

func slidingStartAndEndIndex(pieceType PieceType) (Square, Square) {
	if pieceType == Bishop {
		return 4, 8
	}
	if pieceType == Rook {
		return 0, 4
	}
	return 0, 8
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
			if IgnoreFriendlyPieces {
				kingMoveList.AddMove(Move{From: square, To: to, Flag: NoFlag})
			}
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

func numChecks(pos *Position) int {
	if !SqToEdgeComputed {
		computeSquaresToEdge()
	}
	kingSqaure := pos.GetKingSquare(pos.ColorToMove)

	pos.ColorToMove = pos.ColorToMove.opposite()
	opponentMoves := GetAllPossibleMoves(pos)
	pos.ColorToMove = pos.ColorToMove.opposite()

	count := 0
	previousAttacker := Square(0) // The same attacker can't attack the king twice, this is such that a pawn can't attack the king twice from promotions

	for i := 0; i < opponentMoves.Count; i++ {
		if opponentMoves.Moves[i].To == kingSqaure && previousAttacker != opponentMoves.Moves[i].From {
			count++
			previousAttacker = opponentMoves.Moves[i].From
		}
	}
	return count
}

// IgnoreFriendlyPieces is used to ignore friendly pieces when generating moves. For the Squares under attack BB
var IgnoreFriendlyPieces bool = false

// enemy pieces can also be "Under attack", which makes them not accessable by the king.
// We do this by removing the king from the board, and then generating all possible moves for the enemy pieces.
// There is a bug where enemy pieces should also be include in the mask if they are protected by another piece.
func BBSquaresUnderAttack(pos *Position) Bitboard {
	if !SqToEdgeComputed {
		computeSquaresToEdge()
	}

	pos.Board[pos.GetKingSquare(pos.ColorToMove)] = Piece{NoColor, NoPiece} // Remove the king from the board
	pos.ColorToMove = pos.ColorToMove.opposite()
	var opponentMoves MoveList
	IgnoreFriendlyPieces = true
	for sq := Square(1); sq <= 64; sq++ {
		currPiece := pos.Board[sq]
		if currPiece.Color != pos.ColorToMove {
			continue
		}

		switch currPiece.PieceType {
		case Pawn:
			opponentMoves.AddMoves(genPawnAttackMoves(pos, sq, currPiece))
		case Knight:
			opponentMoves.AddMoves(genKnightMoves(pos, sq, currPiece))
		case Bishop, Queen, Rook:
			opponentMoves.AddMoves(genSlidingMoves(pos, sq, currPiece))
		case King:
			opponentMoves.AddMoves(genKingMoves(pos, sq, currPiece))
		default:
			continue
		}

	}
	IgnoreFriendlyPieces = false
	pos.ColorToMove = pos.ColorToMove.opposite()
	pos.Board[pos.GetKingSquare(pos.ColorToMove)] = Piece{pos.ColorToMove, King} // Add the king back to the board

	var attackedSquares Bitboard

	for i := 0; i < opponentMoves.Count; i++ {
		move := opponentMoves.Moves[i]
		if pos.Board[move.From].PieceType == Pawn {
			if File(move.From) != File(move.To) {
				attackedSquares.SetBit(move.To)
			}
		} else {
			attackedSquares.SetBit(opponentMoves.Moves[i].To)
		}
	}
	return attackedSquares
}

func BBKingAttackedMask(pos *Position) Bitboard {
	if !SqToEdgeComputed {
		computeSquaresToEdge()
	}

	var attackedSquares Bitboard

	pos.ColorToMove = pos.ColorToMove.opposite()
	opponentMoves := GetAllPossibleMoves(pos)
	pos.ColorToMove = pos.ColorToMove.opposite()

	// Figure out the attacker
	var attackerSquare Square = 0

	for i := 0; i < opponentMoves.Count; i++ {
		if pos.Board[opponentMoves.Moves[i].To].PieceType == King {
			attackerSquare = opponentMoves.Moves[i].From
		}
	}

	if attackerSquare == 0 {
		return attackedSquares
	}

	attackedSquares.SetBit(attackerSquare)

	currPiece := pos.Board[attackerSquare]

	if currPiece.PieceType == Pawn || currPiece.PieceType == Knight {
		return attackedSquares
	}

	// Figure out the path to the king
	dirStartIndex, dirEndIndex := slidingStartAndEndIndex(currPiece.PieceType)

	for dirIndex := dirStartIndex; dirIndex < dirEndIndex; dirIndex++ {
		found := false
		for i := int8(1); i <= numSquaresToEdge[attackerSquare][dirIndex]; i++ {
			to := attackerSquare + directionOffsets[dirIndex]*Square(i)

			if to == pos.GetKingSquare(pos.ColorToMove) {
				found = true
				for j := int8(1); j < i; j++ {
					attackedSquares.SetBit(attackerSquare + directionOffsets[dirIndex]*Square(j))
				}
			}

		}
		if found {
			break
		}
	}

	return attackedSquares
}

func BBPinnedSquares(pos *Position) Bitboard {
	var pinnedSquares Bitboard

	if !SqToEdgeComputed {
		computeSquaresToEdge()
	}

	// Go though each enemy piece, and check if it is attacking the king and counting how many pieces are blicking its attack,
	// if there is only one piece blocking the attack, then it is pinned.
	// If there is a enpassent,

	for sq := Square(1); sq <= 64; sq++ {
		currPiece := pos.Board[sq]

		if currPiece.PieceType == NoPiece || currPiece.Color == pos.ColorToMove || currPiece.PieceType == Pawn || currPiece.PieceType == Knight || currPiece.PieceType == King {
			continue
		}

		dirStartIndex, dirEndIndex := slidingStartAndEndIndex(currPiece.PieceType)
		isEnPassantRelevant := pos.EPFile != 0 && dirStartIndex == 0 && ((Rank(sq) == 5 && currPiece.Color == Black) || (Rank(sq) == 4 && currPiece.Color == White)) && Rank(pos.GetKingSquare(pos.ColorToMove)) == Rank(sq)

		for dirIndex := dirStartIndex; dirIndex < dirEndIndex; dirIndex++ {
			piecesGoneThrough := 0
			whiteAndBlackPawnNextToEachOther := 0

			for i := int8(1); i <= numSquaresToEdge[sq][dirIndex]; i++ {
				to := sq + directionOffsets[dirIndex]*Square(i)

				if !isEnPassantRelevant { // Normal sliding moves

					toColor := pos.Board[to].Color
					if toColor == currPiece.Color {
						break
					}

					if pos.GetKingSquare(pos.ColorToMove) == to && piecesGoneThrough == 1 {
						// Set all the previous squares to a pinned square
						for j := int8(1); j < i; j++ {
							pinnedSquares.SetBit(sq + directionOffsets[dirIndex]*Square(j))
						}
						pinnedSquares.SetBit(sq)
						break
					}
					if toColor == currPiece.Color.opposite() {
						if piecesGoneThrough == 1 {
							break
						} else {
							piecesGoneThrough++
						}
					}
				} else { // En passant relevant
					// We have to check that there is exacly a white pawn and a black pawn on the same rank next to each other

					if pos.Board[to].PieceType == Pawn {
						nextSquare := to + directionOffsets[dirIndex]
						if pos.Board[nextSquare].PieceType == Pawn && pos.Board[nextSquare].Color != pos.Board[to].Color {
							whiteAndBlackPawnNextToEachOther = int(nextSquare)
							i++ // Skip the next square, because there is a pawn there, we know
						}
					} else if pos.GetKingSquare(pos.ColorToMove) == to && whiteAndBlackPawnNextToEachOther != 0 {
						// Set all the previous squares to a pinned square
						for j := int8(1); j < i; j++ {
							pinnedSquares.SetBit(sq + directionOffsets[dirIndex]*Square(j))
						}

						if pos.ColorToMove == White {
							pinnedSquares.SetBit(Square(whiteAndBlackPawnNextToEachOther + 8))
						} else {
							pinnedSquares.SetBit(Square(whiteAndBlackPawnNextToEachOther - 8))
						}

						break
					} else if pos.Board[to].PieceType != NoPiece {
						break
					}
				}

			}
		}

	}

	return pinnedSquares
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
