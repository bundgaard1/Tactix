package main

import (
	"tactix/engine"
)

func main() {
	engine.RunCommLoop()
}

// func print_pos(pos *engine.Position) {
// 	fmt.Print(pos.PieceBitboard(engine.Piece{Color: engine.White, PType: engine.King}).StringOnBoard())
// }

// func main() {
// 	pos, _ := engine.FromFEN("k7/1b6/8/8/7R/8/8/4K2R w K - 1 2")

// 	fmt.Print(pos)
// 	print_pos(pos)
// 	m, err := engine.ParseUCIMove(pos, "e1g1")
// 	if err != nil {
// 		panic(err)
// 	}
// 	pos.MakeMove(m)

// 	print_pos(pos)
// 	pos.UndoMove(m)

// 	print_pos(pos)
// }
