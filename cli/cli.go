package cli

import (
	"bufio"
	"fmt"
	"os"
	"tactix/engine"
)

const Banner = `	
\\================================//
||   _____          _   _         ||
||  |_   _|_ _  ___| |_(_)_  __   ||
||    | |/ _  |/ __| __| \ \/ /   ||
||    | | (_| | (__| |_| |>  <    ||
||    |_|\__,_|\___|\__|_/_/\_\   ||
//================================\\
	`

type Cli struct {
	board  engine.Position
	reader *bufio.Reader
}

func NewCli() Cli {
	return Cli{
		board:  engine.FromStandardStartingPosition(),
		reader: bufio.NewReader(os.Stdin),
	}
}

func printPromt() {
	fmt.Print("Tactix> ")
}

func Run() {
	fmt.Println(Banner)

	cli := NewCli()

	for {
		printPromt()
		input, err := cli.reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input")
			continue
		}
		lexer := NewLexer(input)
		command := lexer.GetNextToken()
		cli.handleCommand(command.Literal, &lexer)
	}
}

func (cli *Cli) handleCommand(command string, lexer *Lexer) {
	switch command {
	case "print", "p":
		fmt.Println(cli.board.String())
	case "move", "m":
		cli.handleMoveCommand(lexer)
	case "quit", "q":
		os.Exit(0)
	default:
		fmt.Println("Unknown command : " + command)
	}
}

func (cli *Cli) handleMoveCommand(lexer *Lexer) {

	moveStr := lexer.GetNextToken()

	move, err := engine.ParseUCIMove(&cli.board, moveStr.Literal)
	if err != nil {
		fmt.Println("Invalid move")
		return
	}
	fmt.Println(move.String())

	if !engine.IsMoveValid(&cli.board, move) {
		fmt.Println("Move not legal")
		return
	}

	cli.board.MakeMove(move)
}
