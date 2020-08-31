package main

import (
	"io/ioutil"
	"bufio"
	"os"
	"io"
	"fmt"
)

func main() {
	if len(os.Args) > 1 {
		for _, s := range os.Args[1:] {
			buff, err := ioutil.ReadFile(s)
			if err != nil {
				panic(err)
			}

			tokens := lex(string(buff), s)

			for _, t := range tokens {
				fmt.Println(t.String())
			}
		}
	} else {
		stdin := bufio.NewReader(os.Stdin)
		for {
			l, err := stdin.ReadString('\n')

			tokens := lex(l, "stdin")

			for _, t := range tokens {
				fmt.Println(t.String())
			}

			if err == io.EOF {
				break
			} else if err != nil {
				panic(err)
			}
		}
	}
}
