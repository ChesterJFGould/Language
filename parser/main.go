package main

import (
	"fmt"
	"os"
	"bufio"

	t "../tokens"
)

func main() {
	tokens := make([]t.Token, 0)

	if len(os.Args) > 1 {

	} else {
		stdin := bufio.NewScanner(os.Stdin)

		for stdin.Scan() {
			token, err := t.TokenFromString(stdin.Text())
			if err != nil {
				panic(err)
			}
			tokens = append(tokens, token)
		}
	}

	p := newParser(tokens)

	fmt.Println(p.parseStatement().String())
}
