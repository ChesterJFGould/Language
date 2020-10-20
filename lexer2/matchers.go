package main

import (
	"fmt"
	"unicode"
	"strings"
	"../tokens"
)

var operatorCharset string = "+-*/<>="
var doubleOperatorCharset string = "+-="
var operatorToToken map[string] tokens.TokenType = map[string] tokens.TokenType {
	"+": tokens.Add,
	"++": tokens.Increment,
	"-": tokens.Subtract,
	"--": tokens.Decrement,
	"*": tokens.Multiply,
	"/": tokens.Divide,
	"<": tokens.LessThan,
	">": tokens.GreaterThan,
	"=": tokens.Assignment,
	"==": tokens.EqualTo,
}

var separatorCharset string = ";,(){}[]"

var separatorToToken map[string] tokens.TokenType = map[string] tokens.TokenType {
	";": tokens.Semicolon,
	",": tokens.Comma,
	"(": tokens.OpenBracket,
	")": tokens.CloseBracket,
	"{": tokens.OpenCurlyBracket,
	"}": tokens.CloseCurlyBracket,
	"[": tokens.OpenSquareBracket,
	"]": tokens.CloseSquareBracket,
}

var keywords map[string] tokens.TokenType = map[string] tokens.TokenType {
	"for": tokens.For,
	"if": tokens.If,
	"else": tokens.Else,
	"return": tokens.Return,
	"break": tokens.Break,
	"continue": tokens.Continue,
	"var": tokens.Var,
	"true": tokens.BoolLiteral,
	"false": tokens.BoolLiteral,
}

type stringMatcher struct {}

func (sm stringMatcher) isMatch(r rune) bool {
	return r == '"'
}

func (sm stringMatcher) match(rq *runeQueue) tokens.Token {
	literal := []rune {'"'}
	location := rq.Location
	println(location.String())

	for r, done := rq.next(); r != '"'; r, done = rq.next() {
		if done {
			panic(fmt.Sprintf("Failed to parse string literal at %s, unexpected EOF", rq.Location))
		}

		// Escaped char
		if r == '\\' {
			r, done = rq.next()
			literal = append(literal, '\\', r)
		} else {
			literal = append(literal, r)
		}
	}
	literal = append(literal, '"')
	rq.next()

	return tokens.Token {
		Type: tokens.StringLiteral,
		Literal: string(literal),
		Location: location,
	}
}

type numberMatcher struct {}

func (nm numberMatcher) isMatch(r rune) bool {
	return unicode.IsDigit(r)
}

func (nm numberMatcher) match(rq *runeQueue) tokens.Token {
	r, done := rq.current()
	literal := []rune{r}
	location := rq.Location

	matchDigits := func () {
		for r, done = rq.next(); unicode.IsDigit(r); r, done = rq.next() {
			literal = append(literal, r)
		}
	}
	matchDigits()

	// Float
	if r == '.' {
		literal = append(literal, '.')
		matchDigits()
		if r == 'e' {
			literal = append(literal, 'e')
			matchDigits()
		}

		return tokens.Token {
			Type: tokens.FloatLiteral,
			Literal: string(literal),
			Location: location,
		}
	// Int
	} else {
		if r == 'e' {
			literal = append(literal, 'e')
			matchDigits()
		}

		return tokens.Token {
			Type: tokens.IntLiteral,
			Literal: string(literal),
			Location: location,
		}
	}
}

type operatorMatcher struct {}

func (om operatorMatcher) isMatch(r rune) bool {
	return strings.ContainsRune(operatorCharset, r)
}

func (om operatorMatcher) match(rq *runeQueue) tokens.Token {
	r, _ := rq.current()
	literal := string(r)
	location := rq.Location
	rNext, _ := rq.next()

	if strings.ContainsRune(doubleOperatorCharset, r) && r == rNext {
		literal += string(rNext)
		rq.next()
	}


	return tokens.Token {
		Type: operatorToToken[literal],
		Literal: literal,
		Location: location,
	}
}

type separatorMatcher struct {}

func (sm separatorMatcher) isMatch(r rune) bool {
	return strings.ContainsRune(separatorCharset, r)
}

func (sm separatorMatcher) match(rq *runeQueue) tokens.Token {
	r, _ := rq.current()
	location := rq.Location
	rq.next()

	return tokens.Token {
		Type: separatorToToken[string(r)],
		Literal: string(r),
		Location: location,
	}
}

type keywordIdentifierMatcher struct {}

func (km keywordIdentifierMatcher) isMatch(r rune) bool {
	return !unicode.IsDigit(r) && !strings.ContainsRune(operatorCharset, r) &&
		!strings.ContainsRune(separatorCharset, r) &&
		!unicode.IsSpace(r)
}

func (km keywordIdentifierMatcher) match(rq *runeQueue) tokens.Token {
	r, done := rq.current()
	location := rq.Location
	literal := string(r)

	for r, done = rq.next(); !done && !strings.ContainsRune(operatorCharset, r) &&
					!strings.ContainsRune(separatorCharset, r) &&
					!unicode.IsSpace(r); r, done = rq.next() {
		literal += string(r)
	}

	for k, v := range keywords {
		if k == literal {
			return tokens.Token {
				Type: v,
				Literal: literal,
				Location: location,
			}
		}
	}

	return tokens.Token {
		Type: tokens.Identifier,
		Literal: literal,
		Location: location,
	}
}
