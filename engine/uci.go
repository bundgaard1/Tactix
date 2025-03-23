// Reference UCI interface spec.:
// https://gist.github.com/DOBRO/2592c6dad754ba67e6dcaec8c90165bf#file-uci-protocol-specification-txt

package engine

import (
	"fmt"
	"strconv"
	"strings"
)

type UCI struct {
	options   map[string]string
	pos       *Position
	search    *Search
	open_book *OpeningBook
}

func NewUCI() *UCI {
	return &UCI{
		options:   make(map[string]string),
		pos:       FromStandardStartingPosition(),
		search:    NewSearch(),
		open_book: NewOpeningBook(),
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
		go uci.goCommand(message) // Runs in a goroutine, so we can stop it
	case "position":
		uci.positionCommand(message)
	case "stop":
		uci.stopCommand()
	case "debug":
		uci.debugCommand(message)
	case "ucinewgame":

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

	if uci.open_book.InBook(uci.pos.MoveHistory) {
		move := uci.open_book.GetBookMove(uci.pos)
		uci.replyBestMove(move)
		return
	}

	// Parse the command (options)
	command := strings.TrimPrefix(message, "go ")
	fields := strings.Split(command, " ")

	colorPrefix := "w"
	if uci.pos.ColorToMove == Black {
		colorPrefix = "b"
	}

	timeLeft, increment, moveTime := InfiniteTime, NoValue, NoValue
	movesToGo, maxDepth := NoValue, NoValue
	err := error(nil)

	for i, field := range fields {
		if strings.HasPrefix(field, colorPrefix) {
			if strings.HasSuffix(field, "time") {
				timeLeft, err = strconv.Atoi(fields[i+1])
			}
			if strings.HasSuffix(field, "inc") {
				increment, err = strconv.Atoi(fields[i+1])
			}
		}
		if strings.HasPrefix(field, "movetime") {
			moveTime, err = strconv.Atoi(fields[i+1])
		}
		if strings.HasPrefix(field, "depth") {
			maxDepth, err = strconv.Atoi(fields[i+1])
		}
		if strings.HasPrefix(field, "nodes") {
			movesToGo, err = strconv.Atoi(fields[i+1])
		}
		check(err)
	}

	uci.search.Timer.SetTimeControl(
		int64(timeLeft),
		int64(increment),
		int64(moveTime),
		int64(movesToGo),
		uint8(maxDepth),
	)

	uci.search.Reset()

	uci.search.SetPosition(uci.pos)

	uci.search.Search()
	bestMove := uci.search.BestMove

	uci.replyBestMove(bestMove)
}

func (uci *UCI) replyBestMove(move Move) {
	fmt.Println("bestmove", move.UCIString())
	Log("Best move: " + move.UCIString())
}

func (uci *UCI) stopCommand() {
	uci.search.Stop()
}

func (uci *UCI) positionCommand(message string) {
	msgParts := strings.Split(message, " ")

	if len(msgParts) < 2 {
		fmt.Println("Invalid position command")
		return
	}

	// Moves includes
	movesIncluded, movesIndex := false, 0

	for i, part := range msgParts {
		if part == "moves" {
			movesIncluded, movesIndex = true, i
			break
		}
	}

	// Setup
	if msgParts[1] == "startpos" {
		uci.pos = FromStandardStartingPosition()

	} else if msgParts[1] == "fen" {
		if len(msgParts) < 3 {
			fmt.Println("Invalid position command")
			return
		}
		fenParts := msgParts[2:]
		if movesIncluded {
			fenParts = msgParts[2:movesIndex]
		}
		fen := strings.Join(fenParts, " ")
		pos, err := FromFEN(fen)
		check(err)
		uci.pos = pos
	} else {
		fmt.Println("Invalid position command: unknown position type")
		return
	}
	uci.search.SetPosition(uci.pos)

	// Make the moves stuff
	if movesIncluded {
		for i := movesIndex + 1; i < len(msgParts); i++ {
			moveStr := strings.TrimSpace(msgParts[i])
			move, err := ParseUCIMove(uci.pos, moveStr)
			if err != nil {
				fmt.Println("Invalid move : ", moveStr)
				return
			}
			uci.pos.MakeMove(move)
		}
	}
}

func (uci *UCI) debugCommand(message string) {
	msgParts := strings.Fields(message)
	if len(msgParts) < 2 {
		fmt.Println("Invalid debug command")
		return
	}

	switch msgParts[1] {
	case "on":
		uci.search.Debug = true
	case "off":
		uci.search.Debug = false
	default:
		fmt.Println("Invalid debug command")
	}
}
