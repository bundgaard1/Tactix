package engine

const maxMoveList = 218

type MoveList struct {
	Moves [maxMoveList]Move
	Count uint8
}

func (moveList *MoveList) AddMove(move Move) {
	moveList.Moves[moveList.Count] = move
	moveList.Count++
}
