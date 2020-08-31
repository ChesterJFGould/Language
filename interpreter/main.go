package main

import (
	"bufio"
	"os"

	"../nodes"
)

func main() {
	if len(os.Args) > 1 {

	} else {
		stdin := bufio.NewScanner(os.Stdin)

		if stdin.Scan() {
			s, err := nodes.StatementFromScanner(stdin)
			if err != nil {
				panic(err)
			}

			i := newInterpreter()

			i.interpretStatement(s)
		}
	}
}
