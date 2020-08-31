package main

import (
	"fmt"

	"../tokens"
	"../location"
)

var dfas []*node

func init() {
	dfas = make([]*node, 0)

	for _, dfa := range dfaTemplates {
		nodes := parseNodeTemplates(dfa)

		for _, n := range nodes {
			if n.isStarter {
				dfas = append(dfas, n)
			}
		}
	}
}

type runeQueue struct {
	i int
	queue []rune
	location.Location
}

func (rq *runeQueue) next() (rune, bool) {
	if rq.i < len(rq.queue) {
		if rq.queue[rq.i] == '\n' {
			rq.Location.Line += 1
			rq.Location.Column = 1
		} else {
			rq.Location.Column += 1
		}

		defer func() {
			rq.i++
		}()

		return rq.queue[rq.i], false
	} else {
		return 0, true
	}
}

func (rq *runeQueue) peek() (rune, bool) {
	if rq.i < len(rq.queue) {
		return rq.queue[rq.i], false
	} else {
		return 0, true
	}
}

func runeQueueFromString(s string) runeQueue {
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

type node struct {
	isStarter bool
	isMatch func(r rune) bool
	connections []*node
	canEmit bool
	emittedType tokens.TokenType
}

func (n *node) parse(rq *runeQueue) tokens.Token {
	location := rq.Location

	r, done := rq.next()
	matched := string(r)

	if done && n.canEmit {
		return tokens.Token {
			Type: n.emittedType,
			Literal: matched,
			Location: location,
		}
	} else if done {
		panic(fmt.Sprintf("lexer failed to match %s", location))
	}

	for {
		r, done := rq.peek()
		if done && n.canEmit {
			return tokens.Token {
				Type: n.emittedType,
				Literal: matched,
				Location: location,
			}
		} else if done {
			panic(fmt.Sprintf("lexer failed to match %s", location))
		}

		foundMatch := false

		for _, c := range n.connections {
			if c.isMatch(r) {
				foundMatch = true
				matched += string(r)
				n = c
				break
			}
		}

		if foundMatch {
			rq.next()
			continue
		} else if n.canEmit {
			return tokens.Token {
				Type: n.emittedType,
				Literal: matched,
				Location: location,
			}
		} else {
			panic(fmt.Sprintf("lexer failed to match %s", location))
		}
	}
}

type nodeTemplate struct {
	isStarter bool
	isMatch func(r rune) bool
	connections []int
	canEmit bool
	emittedType tokens.TokenType
}

func parseNodeTemplates(nts []nodeTemplate) []*node {
	nodes := make([]*node, len(nts))

	for i, n := range nts {
		nodes[i] = &node{}
		nodes[i].isStarter = n.isStarter
		nodes[i].isMatch = n.isMatch
		nodes[i].connections = make([]*node, len(n.connections))
		nodes[i].canEmit = n.canEmit
		nodes[i].emittedType = n.emittedType
	}

	for i, nt := range nts {
		for j, n := range nt.connections {
			nodes[i].connections[j] = nodes[n]
		}
	}

	return nodes
}

func lex(s string, fileName string) []tokens.Token {
	rq := runeQueueFromString(s)
	rq.Location.File = fileName

	tokens := make([]tokens.Token, 0)

	for r, done := rq.peek(); !done; r, done = rq.peek() {
		matched := false

		for _, n := range dfas {
			if n.isMatch(r) {
				matched = true
				tokens = append(tokens, n.parse(&rq))
				break
			}
		}

		if !matched {
			rq.next()
		}
	}

	return tokens
}
