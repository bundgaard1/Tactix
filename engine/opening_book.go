package engine

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type obNode struct {
	uciMove  string
	children []*obNode
}

func newNode(move string) *obNode {
	return &obNode{
		uciMove: move,
	}
}

func (bn *obNode) contains(move string) (bool, *obNode) {
	for _, child := range bn.children {
		if child.uciMove == move {
			return true, child
		}
	}
	return false, nil
}

func (bn *obNode) addChild(move string) *obNode {
	if exists, node := bn.contains(move); exists {
		return node
	}
	nn := newNode(move)

	bn.children = append(bn.children, nn)

	return nn
}

type OpeningBook struct {
	Root *obNode
	Curr *obNode
}

func NewOpeningBook() *OpeningBook {
	ob := OpeningBook{}
	ob.Root = newNode("ROOT")
	ob.Curr = ob.Root

	dat, err := os.Open("opening_book/openings.txt")
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
	builder.WriteString("\n")
	indent += 1
	fmt.Println(indent)
	for i, child := range bn.children {
		builder.WriteString(strings.Repeat(" ", indent))
		builder.WriteString(fmt.Sprint(i+1, ": "))
		builder.WriteString(child.String())
	}
	indent -= 1
	return builder.String()
}
