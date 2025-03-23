package engine

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	EngineName   = "Tactix 0.1"
	EngineAuthor = "Mathias Jørgensen"

	Banner = `	
	\\================================//
	||   _____          _   _         ||
	||  |_   _|_ _  ___| |_(_)_  __   ||
	||    | |/ _  |/ __| __| \ \/ /   ||
	||    | | (_| | (__| |_| |>  <    ||
	||    |_|\__,_|\___|\__|_/_/\_\   ||
	//================================\\
	`

	HelpMessage = `	Commands:
	uci - Start th UCI protocol
	d/print - Display the current board
	move <move> - Make a move
	perft <depth> - Run perft to a certain depth
	position <fen> - Set the board to a fen string
	help - Print this help message
	quit - Exit the program
`
)

type CommChannel struct {
	reader *bufio.Reader
	uci    *UCI
}

func NewComms() *CommChannel {
	com := CommChannel{
		reader: bufio.NewReader(os.Stdin),
	}
	com.uci = NewUCI()
	return &com
}

func RunCommLoop() {
	fmt.Println(Banner)

	InitLogging()
	comms := NewComms()

	for {
		message, err := comms.reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input")
			continue
		}
		message = strings.Trim(message, "\n")
		if message == "quit" {
			break
		}
		Log("Received command: " + message)
		comms.handleCommand(message)
	}
}

func (comm *CommChannel) handleCommand(message string) {
	fields := strings.Fields(message)
	if len(fields) == 0 {
		return
	}

	switch fields[0] {
	// UCI commands
	case "uci", "isready", "setoption", "register", "ucinewgame", "go", "position", "stop", "ponderhit":
		comm.uci.handleUCICommand(message)
	// Custom commands
	case "d", "print":
		fmt.Println(comm.uci.pos.String())
	case "eval":
		fmt.Println(Evaluate(comm.uci.pos))
	case "move", "m":
		comm.moveCommand(message)
	case "perft":
		comm.perftCommand(message)
	case "moves":
		moves := LegalMoves(comm.uci.pos)
		fmt.Println(moves.String())
	case "help", "h":
		comm.helpCommand()
	case "bench":
		PerftWithBenchmark()
	default:
		fmt.Println("Unknown command : " + fields[0])
		comm.handleCommand(strings.Join(fields[1:], " "))
	}
}

func (comm *CommChannel) moveCommand(message string) {
	msgParts := strings.Fields(message)
	if len(msgParts) < 2 {
		fmt.Println("Invalid move command")
		return
	}

	move, err := ParseUCIMove(comm.uci.pos, msgParts[1])
	if err != nil {
		fmt.Println("Invalid move")
		return
	}

	fmt.Println(move.String())

	if !IsMoveValid(comm.uci.pos, move) {
		fmt.Println("Move not legal")
		return
	}

	comm.uci.pos.MakeMove(move)
}

var PerftCommandHelp = "perft <depth>"

func (comm *CommChannel) perftCommand(message string) {
	msgParts := strings.Fields(message)
	if len(msgParts) < 2 {
		fmt.Println("No depth Provided:", PerftCommandHelp)
		return
	}

	depth, err := strconv.Atoi(msgParts[1])
	if err != nil {
		fmt.Println("Invalid Depth:", PerftCommandHelp)
		return
	}

	summary, nodes := PerftDivided(comm.uci.pos, depth)

	fmt.Println(summary)
	fmt.Println("Total nodes:", nodes)
}

func (comm *CommChannel) helpCommand() {
	fmt.Print(HelpMessage)
}
