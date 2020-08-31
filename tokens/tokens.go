package tokens

import (
	"fmt"
	"strings"

	"../location"
)

var tokenTypeToString map[TokenType] string = map[TokenType] string {
	IntLiteral: "IntLiteral",
	FloatLiteral: "FloatLiteral",
	StringLiteral: "StringLiteral",
	Add: "Add",
	Increment: "Increment",
	Subtract: "Subtract",
	Decrement: "Decrement",
	Multiply: "Multiply",
	Divide: "Divide",
	LessThan: "LessThan",
	GreaterThan: "GreaterThan",
	EqualTo: "EqualTo",
	Identifier: "Identifier",
	For: "For",
	If: "If",
	Else: "Else",
	Return: "Return",
	Break: "Break",
	Continue: "Continue",
	Var: "Var",
	Semicolon: "Semicolon",
	Comma: "Comma",
	OpenBracket: "OpenBracket",
	CloseBracket: "CloseBracket",
	OpenCurlyBracket: "OpenCurlyBracket",
	CloseCurlyBracket: "CloseCurlyBracket",
	OpenSquareBracket: "OpenSquareBracket",
	CloseSquareBracket: "CloseSquareBracket",
	Assignment: "Assignment",
	EOF: "EOF",
}

// Generated in init as reverse of tokenTypeToString
var stringToTokenType map[string]TokenType = map[string] TokenType {}

func init() {
	// Generate stringToTokenType
	for k, v := range tokenTypeToString {
		stringToTokenType[v] = k
	}
}

type TokenType int

const (
	IntLiteral TokenType = iota + 1
	FloatLiteral
	StringLiteral
	Add
	Increment
	Subtract
	Decrement
	Multiply
	Divide
	LessThan
	GreaterThan
	EqualTo
	Identifier
	For
	If
	Else
	Return
	Break
	Continue
	Var
	Semicolon
	Comma
	OpenBracket
	CloseBracket
	OpenCurlyBracket
	CloseCurlyBracket
	OpenSquareBracket
	CloseSquareBracket
	Assignment
	EOF
)

func (t TokenType) String() string {
	return tokenTypeToString[t]
}

func TokenTypeFromString(s string) (TokenType, error) {
	if t, ok := stringToTokenType[s]; ok {
		return t, nil
	} else {
		return 0, fmt.Errorf("can't parse %q into TokenType", s)
	}
}


type Token struct {
	Type TokenType
	Literal string
	location.Location
}

func (t Token) String() string {
	// Replace all spaces in literal with zero width spaces.
	// Although this kind of feels like a hack it's kind of
	// what zero width spaces are meant for.
	// t.Literal = strings.ReplaceAll(t.Literal, " ", string(rune(8203)))
	return fmt.Sprintf("%s %s %s", t.Type, t.Literal, t.Location)
}

func TokenFromString(s string) (Token, error) {
	t := Token{}

	vals := strings.Split(s, " ")

	var err error

	t.Type, err = TokenTypeFromString(vals[0])
	if err != nil {
		return t, err
	}

	// Replace all zero width spaces with spaces
	//t.Literal = strings.ReplaceAll(vals[1], string(rune(8203)), " ")

	t.Literal = strings.Join(vals[1:len(vals)-3], " ")

	t.Location, err = location.LocationFromString(strings.Join(vals[len(vals)-3:], " "))
	if err != nil {
		return t, err
	}

	return t, nil
}
