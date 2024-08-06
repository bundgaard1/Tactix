package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
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
	writer *bufio.Writer
}

func NewCli() Cli {
	return Cli{
		board:  engine.FromStandardStartingPosition(),
		reader: bufio.NewReader(os.Stdin),
		writer: bufio.NewWriter(os.Stdout),
	}
}

func printPromt() {
	fmt.Print("")
}

func (cli *Cli) toOutput(s string) {
	cli.writer.Write([]byte(s))
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

		if lexer.IsAtEnd() {
			continue
		}

		cli.handleCommand(command.Literal, &lexer)
	}
}

func (cli *Cli) handleCommand(command string, lexer *Lexer) {

	switch command {
	case "quit":
		os.Exit(0)
	case "d":
		fmt.Println(cli.board.String())
	case "perft":
		cli.perftCommand(lexer)
	case "bench":
		cli.benchCommand(lexer)
	case "move", "m":
		cli.moveCommand(lexer)
	case "position", "pos":
		cli.positionCommand(lexer)
	case "help", "h":
		cli.helpCommand()
	default:
		fmt.Println("Unknown command : " + command)
	}
}

func (cli *Cli) moveCommand(lexer *Lexer) {

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
	fmt.Println(cli.board.String())

}

func (cli *Cli) perftCommand(lexer *Lexer) {
	depthStr := lexer.GetNextToken()

	depth, err := strconv.Atoi(depthStr.Literal)

	if err != nil {
		fmt.Println("Invalid depth")
		return
	}

	summary, nodes := engine.PerftDivided(&cli.board, depth)

	fmt.Println(summary)
	fmt.Println("Total nodes: ", nodes)
}

func (cli *Cli) positionCommand(lexer *Lexer) {
	fenStr := lexer.GetRestOfInput()
	cli.board = engine.FromFEN(fenStr)
	fmt.Println(cli.board.String())
}

func (cli *Cli) benchCommand(lexer *Lexer) {
	fmt.Println("Bench command")
}

func (cli *Cli) helpCommand() {
	fmt.Println("Commands:")
	fmt.Println("	d - Display the current board")
	fmt.Println("	move <move> - Make a move")
	fmt.Println("	perft <depth> - Run perft to a certain depth")
	fmt.Println("	position <fen> - Set the board to a fen string")
	fmt.Println("	help - Print this help message")
	fmt.Println("	quit - Exit the program")
}
