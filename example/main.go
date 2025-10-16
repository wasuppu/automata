package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/wasuppu/automata"
)

func main() {
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		fmt.Fprintf(os.Stderr, "usage: automata -E <pattern>\n")
		os.Exit(2)
	}

	pattern := os.Args[2]
	nfa := automata.Interp(pattern)

	inputReader := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stdout, "> ")
		line, _, err := inputReader.ReadLine()

		if errors.Is(err, io.EOF) {
			break
		}

		if !nfa.Matches(string(line)) {
			fmt.Printf("%s\t\t=> Not matched.\n", string(line))
		} else {
			fmt.Printf("%s\t\t=> Matched.\n", string(line))
		}
	}
}
