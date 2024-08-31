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

type Communication struct {
	pos    *Position
	reader *bufio.Reader
	uci    *UCI
}

func NewComms() *Communication {
	com := Communication{
		pos:    FromStandardStartingPosition(),
		reader: bufio.NewReader(os.Stdin),
	}
	com.uci = NewUCI(com.pos)
	return &com
}

func RunCommLoop() {
	fmt.Println(Banner)

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

		comms.handleCommand(message)
	}
}

func (comm *Communication) handleCommand(message string) {
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
		fmt.Println(comm.pos.String())
	case "eval":
		fmt.Println(Evaluate(comm.pos))
	case "move", "m":
		comm.moveCommand(message)
	case "perft":
		comm.perftCommand(message)
	case "moves":
		moves := LegalMoves(comm.pos)
		fmt.Println(moves.String())
	case "help", "h":
		comm.helpCommand()
	default:
		// If the command is unknown, it should try to parse the remaning part as a uci command,
		// not just skip the rest
		fmt.Println("Unknown command : " + fields[0])
	}
}

func (comm *Communication) moveCommand(message string) {
	msgParts := strings.Fields(message)
	if len(msgParts) < 2 {
		fmt.Println("Invalid move command")
		return
	}

	move, err := ParseUCIMove(comm.pos, msgParts[1])
	if err != nil {
		fmt.Println("Invalid move")
		return
	}

	fmt.Println(move.String())

	if !IsMoveValid(comm.pos, move) {
		fmt.Println("Move not legal")
		return
	}

	comm.pos.MakeMove(move)
}

func (comm *Communication) perftCommand(message string) {
	msgParts := strings.Fields(message)

	depth, err := strconv.Atoi(msgParts[1])
	if err != nil {
		fmt.Println("Invalid depth")
		return
	}

	summary, nodes := PerftDivided(comm.pos, depth)

	fmt.Println(summary)
	fmt.Println("Total nodes: ", nodes)
}

func (comm *Communication) helpCommand() {
	fmt.Print(HelpMessage)
}
