package engine

import (
	"fmt"
	"strings"
)

type MoveList []Move

const initialMoveListSize = 32

func NewMoveList() *MoveList {
	moveList := make(MoveList, 0, initialMoveListSize)
	return &moveList
}

func (ml *MoveList) Append(moves ...Move) {
	l := len(*ml)
	if l+len(moves) > cap(*ml) {
		newList := make(MoveList, (l+len(moves))*2)
		copy(newList, *ml)
		*ml = newList // Update the pointer to the new list
	}
	*ml = (*ml)[0 : l+len(moves)]
	copy((*ml)[l:], moves)
}

func (ml *MoveList) Get(index int) *Move {
	return &(*ml)[index]
}

func (ml *MoveList) AppendList(append *MoveList) {
	ml.Append((*append)...)
}

func (ml *MoveList) String() string {
	var builder strings.Builder
	builder.WriteString("Moves: ")
	builder.WriteString(fmt.Sprintf("%d\n", len(*ml)))

	for i := 0; i < len(*ml); i++ {
		move := (*ml)[i]
		builder.WriteString(move.UCIString())
		builder.WriteString(" ")
	}

	return builder.String()
}
