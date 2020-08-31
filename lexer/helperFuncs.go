package main

import (
	"strings"
)

func matchRune(matcherRune rune) func(r rune) bool {
	return func(r rune) bool {
		return r == matcherRune
	}
}

func matchRuneInverse(matcherRune rune) func(r rune) bool {
	return func(r rune) bool {
    		return r != matcherRune
	}
}

func matchCharset(charset string) func(r rune) bool {
	return func(r rune) bool {
		return strings.ContainsRune(charset, r)
	}
}

func matchCharsetInverse(charset string) func(r rune) bool {
	return func(r rune) bool {
		return !strings.ContainsRune(charset, r)
	}
}

func matchAny(r rune) bool {
	return true
}
