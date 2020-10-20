package main

import (
	"../tokens"
	"../location"
	"os"
	"io/ioutil"
	"fmt"
)

type runeQueue struct {
	i int
	queue []rune
	location.Location
}

func (rq *runeQueue) next() (rune, bool) {
	if rq.i < len(rq.queue) - 1 {
		if rq.queue[rq.i] == '\n' {
			rq.Location.Line += 1
			rq.Location.Column = 1
		} else {
			rq.Location.Column += 1
		}

		rq.i++

		return rq.queue[rq.i], false
	} else if rq.i == len(rq.queue) - 1 {
		rq.i++
		return 0, true
	} else {
		return 0, true
	}
}

func (rq *runeQueue) current() (rune, bool) {
	if rq.i < len(rq.queue) {
		return rq.queue[rq.i], false
	} else {
		return 0, true
	}
}

func (rq runeQueue) fromString(s string) runeQueue {
	return runeQueue {
		i: 0,
		queue: []rune(s),
		Location: location.Location {
			File: "",
			Line: 1,
			Column: 1,
		},
	}
}

type matcher interface {
	isMatch(first rune) bool
	match(rq *runeQueue) tokens.Token
}

var matchers []matcher

func init() {
	matchers = []matcher {
		stringMatcher{},
		numberMatcher{},
		operatorMatcher{},
		separatorMatcher{},
		keywordIdentifierMatcher{},
	}
}

func main() {
	if len(os.Args) > 1 {
		for _, fileName := range os.Args[1:] {
			buff, err := ioutil.ReadFile(fileName)
			if err != nil {
				panic(err)
			}

			tokens := lex(string(buff), fileName)

			for _, t := range tokens {
				fmt.Println(t.String())
			}
		}
	} else {
		// Read from stdin
	}
}

func lex(s string, fileName string) []tokens.Token {
	rq := runeQueue{}.fromString(s)
	rq.Location.File = fileName

	tokens := make([]tokens.Token, 0)

	for r, done := rq.current(); !done; r, done = rq.current() {
		matched := false
		for _, m := range matchers {
			if m.isMatch(r) {
				matched = true
				tokens = append(tokens, m.match(&rq))
				break
			}
		}

		if !matched {
			rq.next()
		}
	}

	return tokens
}
