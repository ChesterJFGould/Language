package main

import (
	"unicode"
	"strings"

	"../tokens"
)

var operatorCharset = "-+/*<>"
var separatorCharset = ";,(){}[]"
var assignmentCharset = "="

var dfaTemplates [][]nodeTemplate = [][]nodeTemplate {
	// String Literal
	{
		{true, matchRune('"'), []int{1, 3, 2}, false, 0},
		{false, matchRune('\\'), []int{0, 2}, false, 0},
		{false, matchAny, []int{1, 3, 2}, false, 0},
		{false, matchRune('"'), []int{}, true, tokens.StringLiteral},
	},
	// Number Literals
	{
		{true, unicode.IsDigit, []int{0, 1, 3}, true, tokens.IntLiteral},
		{false, matchCharset("eE"), []int{2}, false, 0},
		{false, unicode.IsDigit, []int{2}, true, tokens.IntLiteral},
		{false, matchRune('.'), []int{4}, false, 0},
		{false, unicode.IsDigit, []int{4, 5}, true, tokens.FloatLiteral},
		{false, matchCharset("eE"), []int{6}, false, 0},
		{false, unicode.IsDigit, []int{6}, true, tokens.FloatLiteral},
	},
	// Operators
	{
		{true, matchRune('+'), []int{1}, true, tokens.Add},
		{false, matchRune('+'), []int{}, true, tokens.Increment},
		{true, matchRune('-'), []int{3}, true, tokens.Subtract},
		{false, matchRune('-'), []int{}, true, tokens.Decrement},
		{true, matchRune('*'), []int{}, true, tokens.Multiply},
		{true, matchRune('/'), []int{}, true, tokens.Divide},
		{true, matchRune('<'), []int{}, true, tokens.LessThan},
		{true, matchRune('>'), []int{}, true, tokens.GreaterThan},
		{true, matchRune('='), []int{9}, true, tokens.Assignment},
		{true, matchRune('='), []int{}, true, tokens.EqualTo},
	},
	// Separators
	{
		{true, matchRune(';'), []int{}, true, tokens.Semicolon},
		{true, matchRune(','), []int{}, true, tokens.Comma},
		{true, matchRune('('), []int{}, true, tokens.OpenBracket},
		{true, matchRune(')'), []int{}, true, tokens.CloseBracket},
		{true, matchRune('{'), []int{}, true, tokens.OpenCurlyBracket},
		{true, matchRune('}'), []int{}, true, tokens.CloseCurlyBracket},
		{true, matchRune('['), []int{}, true, tokens.OpenSquareBracket},
		{true, matchRune(']'), []int{}, true, tokens.CloseSquareBracket},
	},
	// Keywords and Identifiers
	generateKeywordsAndIdentifiers(
		operatorCharset+separatorCharset+assignmentCharset,
		[]string{"for", "if", "else", "return", "break", "continue", "var"},
		[]tokens.TokenType{tokens.For, tokens.If, tokens.Else, tokens.Return, tokens.Break, tokens.Continue, tokens.Var},
	),
}

func generateKeywordsAndIdentifiers(invalidCharset string, keywords []string, tokenTypes []tokens.TokenType) []nodeTemplate {
	nodes := make([]nodeTemplate, 0)

	for i, s := range keywords {
		nodes = append(nodes, nodeTemplate{true, matchRune([]rune(s)[0]),
		               []int{len(nodes)+1}, true, tokens.Identifier})

		for _, r := range s[1:len(s)-1] {
			nodes = append(nodes, nodeTemplate{false, matchRune(r),
			               []int{len(nodes)+1}, true, tokens.Identifier})
		}

		nodes = append(nodes, nodeTemplate{false, matchRune([]rune(s)[len(s)-1]),
				[]int{}, true, tokenTypes[i]})
	}

	for i, n := range nodes {
		nodes[i].connections = append(n.connections, len(nodes)+1)
	}

	nodes = append(nodes, nodeTemplate{true, func(r rune) bool {
		return !unicode.IsDigit(r) && !strings.ContainsRune(invalidCharset, r) && !unicode.IsSpace(r) && unicode.IsGraphic(r)
	}, []int{len(nodes)+1}, true, tokens.Identifier})

	nodes = append(nodes, nodeTemplate{false, func(r rune) bool {
		return !strings.ContainsRune(invalidCharset, r) && !unicode.IsSpace(r) && unicode.IsGraphic(r)
	}, []int{len(nodes)}, true, tokens.Identifier})

	return nodes
}
