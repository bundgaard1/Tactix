package engine

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
)

type obNode struct {
	uciMove  string
	children []*obNode
}

func newNode(move string) *obNode {
	return &obNode{
		uciMove:  move,
		children: make([]*obNode, 0),
	}
}

func (bn *obNode) Contains(move string) (bool, *obNode) {
	for _, child := range bn.children {
		if child.uciMove == move {
			return true, child
		}
	}
	return false, nil
}

func (bn *obNode) addChild(move string) *obNode {
	if exists, node := bn.Contains(move); exists {
		return node
	}
	nn := newNode(move)

	bn.children = append(bn.children, nn)

	return nn
}

type OpeningBook struct {
	Root *obNode
}

func NewOpeningBook() *OpeningBook {
	ob := OpeningBook{}
	ob.Root = newNode("ROOT")

	dat, err := os.Open("resources/book_openings.txt")
	check(err)
	reader := bufio.NewReader(dat)

	reader.ReadString('\n')
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		check(err)

		moves := strings.Split(line, " ")
		curr := ob.Root
		for _, move := range moves {
			nn := curr.addChild(move)
			curr = nn
		}
	}

	return &ob
}

// Can be used to get a move from the opening book, only call if InBook returns true
func (ob *OpeningBook) GetBookMove(pos *Position) Move {
	curr := ob.Root
	moves := pos.MoveHistory
	for _, move := range *moves {
		if exists, node := curr.Contains(move.UCIString()); exists {
			curr = node
		} else {
			break
		}
	}

	for i, child := range curr.children {
		fmt.Printf("%d: %s\n", i, child.uciMove)
	}

	mIdx := rand.Int() % (len(curr.children))
	m, err := ParseUCIMove(pos, curr.children[mIdx].uciMove)
	if err != nil {
		fmt.Printf("%+v  %+v \n", mIdx, len(curr.children))
		panic("invalid move, tried to parse: " + curr.children[mIdx].uciMove)
	}

	return m
}

func (ob *OpeningBook) InBook(moves *MoveList) bool {
	curr := ob.Root
	for _, move := range *moves {
		if exists, node := curr.Contains(move.UCIString()); exists {
			curr = node
		} else {
			return false
		}
	}
	return !(len(curr.children) == 0)
}

func (ob *OpeningBook) String() string {
	return ob.Root.String()
}

func (bn *obNode) ChildrenString() string {
	var builder strings.Builder
	for _, child := range bn.children {
		builder.WriteString(child.uciMove)
		builder.WriteString(" ")

	}
	return builder.String()
}

var indent int = 0

func (bn *obNode) String() string {
	var builder strings.Builder

	builder.WriteString(bn.uciMove)
	indent += 1
	builder.WriteString(fmt.Sprint(" "))
	if len(bn.children) == 1 {
		builder.WriteString(bn.children[0].String())
	} else {
		for i, child := range bn.children {
			builder.WriteString("\n")
			builder.WriteString(strings.Repeat(" ", indent*5))
			builder.WriteString(fmt.Sprintf("%d. ", i+1))
			builder.WriteString(child.String())
		}
	}
	indent -= 1
	return builder.String()
}
