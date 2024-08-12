package main

import (
	"fmt"
	"tactix/cli"
	"tactix/engine"
	"time"
)

func main() {
	cli.Run()

	now := time.Now()

	engine.NewOpeningBook()
	after := time.Now()

	fmt.Println(after.Sub(now))
}
