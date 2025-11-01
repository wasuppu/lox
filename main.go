package main

import (
	"fmt"
	lox "lox/treewalk"
	"os"
)

func main() {
	l := lox.New()

	if len(os.Args) > 2 {
		fmt.Println("Usage: lox [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		l.RunFile(os.Args[1])
	} else {
		l.RunPrompt(os.Stdin, os.Stdout)
	}
}
