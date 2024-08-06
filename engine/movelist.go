package engine

const maxMoveList = 218

type MoveList struct {
	Moves [maxMoveList]Move
	Count int
}

func (moveList *MoveList) AddMove(move Move) {
	if moveList.Count < maxMoveList {
		moveList.Moves[moveList.Count] = move
		moveList.Count++
	} else {
		panic("MoveList.Count Limit reached, somehow")
	}
}

func (moveList *MoveList) AddMoves(appendList MoveList) {
	for i := 0; i < appendList.Count; i++ {
		moveList.AddMove(appendList.Moves[i])
	}
}

func (MoveList *MoveList) Remove(index int) {
	if index < MoveList.Count {
		MoveList.Moves[index] = MoveList.Moves[MoveList.Count-1]
		MoveList.Count--
	}
}
