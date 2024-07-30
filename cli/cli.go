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
		if lexer.IsAtEnd() {
			continue
		}
		cli.handleCommand(command.Literal, &lexer)
	}
}

func (cli *Cli) handleCommand(command string, lexer *Lexer) {
	switch command {
	case "print", "p":
		fmt.Println(cli.board.String())
	case "perft":
		cli.handlePerftCommand(lexer)
	case "move", "m":
		cli.handleMoveCommand(lexer)
	case "fen":
		fmt.Println(engine.FEN(&cli.board))
	case "position", "pos":
		cli.handlePositionCommand(lexer)
	case "help", "h":
		cli.handleHelpCommand()
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
	fmt.Println(cli.board.String())

}

func (cli *Cli) handlePerftCommand(lexer *Lexer) {
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

func (cli *Cli) handlePositionCommand(lexer *Lexer) {
	fenStr := lexer.GetRestOfInput()
	cli.board = engine.FromFEN(fenStr)
	fmt.Println(cli.board.String())
}

func (cli *Cli) handleHelpCommand() {
	fmt.Println("Commands:")
	fmt.Println("	print - Print the current board")
	fmt.Println("	move <move> - Make a move")
	fmt.Println("	perft <depth> - Run perft to a certain depth")
	fmt.Println("	fen - Print the board from a fen string")
	fmt.Println("	position <fen> - Set the board to a fen string")
	fmt.Println("	help - Print this help message")
	fmt.Println("	quit - Exit the program")
}
