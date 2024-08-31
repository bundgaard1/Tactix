// Reference UCI interface spec.:
// https://gist.github.com/DOBRO/2592c6dad754ba67e6dcaec8c90165bf#file-uci-protocol-specification-txt

package engine

import (
	"fmt"
	"strings"
)

type UCI struct {
	options   map[string]string
	pos       *Position
	open_book *OpeningBook

	// Options

	Debug bool
}

func NewUCI(pos *Position) *UCI {
	return &UCI{
		options:   make(map[string]string),
		pos:       pos,
		open_book: NewOpeningBook(),

		Debug: false,
	}
}

func (uci *UCI) handleUCICommand(message string) {
	fields := strings.Fields(message)
	if len(fields) == 0 {
		return
	}

	switch fields[0] {
	case "uci":
		uci.respondUCI()
	case "isready":
		fmt.Print("readyok\n")
	case "go":
		uci.goCommand(message)
	case "position":
		uci.positionCommand(message)
	default:
		fmt.Print("UCI command not implemented\n")
	}
}

func (uci *UCI) respondUCI() {

	// Engine Identification
	fmt.Print("id name ", EngineName, "\n")
	fmt.Print("id author ", EngineAuthor, "\n")

	// Engine Options
	fmt.Print("option name OwnBook type check default true\n")

	fmt.Print("uciok\n")
}

func (uci *UCI) goCommand(message string) {
	uci.pos.PrintHistory()
	if uci.open_book.InBook(uci.pos.MoveHistory) {
		move := uci.open_book.GetBookMove(uci.pos)
		fmt.Println("BestMove: ", move.UCIString())
		return
	}

	search := NewSearch(uci.pos)

	go search.Search()
	bestMove := search.BestMove

	fmt.Println("BestMove: ", bestMove.UCIString())
}

func (uci *UCI) positionCommand(message string) {
	msgParts := strings.Split(message, " ")
	// TODO : make it handle "moves ..""
	if len(msgParts) < 2 {
		fmt.Println("Invalid position command")
		return
	}

	// Setu
	if msgParts[1] == "startpos" {
		uci.pos = FromStandardStartingPosition()

	} else if msgParts[1] == "fen" {

		fen := strings.Join(msgParts[2:], " ")
		pos, err := FromFEN(fen)
		check(err)
		uci.pos = pos
	}

}
