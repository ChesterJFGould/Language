package main

import (
	"strings"
	"fmt"
	"strconv"
	"math"

	"../tokens"
	"../nodes"
)

////////////////
// Precedence //
////////////////

type Precedence int

const (
	_ Precedence = iota
	LOWEST
	EQUALITY
	SUM
	PRODUCT
	UNARY
	CALL
)

var tokenTypeToPrecedence map[tokens.TokenType] Precedence = map[tokens.TokenType] Precedence {
	tokens.LessThan: EQUALITY,
	tokens.GreaterThan: EQUALITY,
	tokens.EqualTo: EQUALITY,
	tokens.Add: SUM,
	tokens.Subtract: SUM,
	tokens.Multiply: PRODUCT,
	tokens.Divide: PRODUCT,
	tokens.Increment: UNARY,
	tokens.Decrement: UNARY,
	tokens.OpenBracket: CALL,
	tokens.OpenSquareBracket: CALL,
}

func getPrecedence(t tokens.TokenType) Precedence {
	if p, exists := tokenTypeToPrecedence[t]; exists {
		return p
	} else {
		return LOWEST
	}
}

////////////
// parser //
////////////

type parser struct {
	tokens []tokens.Token
	tokenPosition int
	prefixParsers map[tokens.TokenType] func() nodes.Expression
	infixParsers map[tokens.TokenType] func(left nodes.Expression) nodes.Expression
	statementParsers map[tokens.TokenType] func() nodes.Statement
}

func newParser(t []tokens.Token) *parser {
	p := &parser{}

	p.tokens = t
	p.tokenPosition = 0

	p.prefixParsers = make(map[tokens.TokenType] func() nodes.Expression)
	p.infixParsers = make(map[tokens.TokenType] func(left nodes.Expression) nodes.Expression)
	p.statementParsers = make(map[tokens.TokenType] func() nodes.Statement)

	p.prefixParsers[tokens.IntLiteral] = p.parseIntLiteral
	p.prefixParsers[tokens.FloatLiteral] = p.parseFloatLiteral
	p.prefixParsers[tokens.StringLiteral] = p.parseStringLiteral
	p.prefixParsers[tokens.Identifier] = p.parseIdentifier
	p.prefixParsers[tokens.Increment] = p.parsePrefixOperator
	p.prefixParsers[tokens.Decrement] = p.parsePrefixOperator
	p.prefixParsers[tokens.Subtract] = p.parsePrefixOperator
	p.prefixParsers[tokens.OpenBracket] = p.parseSubExpression

	p.infixParsers[tokens.LessThan] = p.parseOperator
	p.infixParsers[tokens.GreaterThan] = p.parseOperator
	p.infixParsers[tokens.EqualTo] = p.parseOperator
	p.infixParsers[tokens.Add] = p.parseOperator
	p.infixParsers[tokens.Subtract] = p.parseOperator
	p.infixParsers[tokens.Multiply] = p.parseOperator
	p.infixParsers[tokens.Divide] = p.parseOperator
	p.infixParsers[tokens.OpenBracket] = p.parseCall
	p.infixParsers[tokens.OpenSquareBracket] = p.parseIndex

	p.statementParsers[tokens.If] = p.parseIf
	p.statementParsers[tokens.For] = p.parseFor
	p.statementParsers[tokens.Var] = p.parseAssignment
	p.statementParsers[tokens.OpenCurlyBracket] = p.parseScope

	return p
}

func (p *parser) nextToken() {
	if p.tokenPosition < len(p.tokens)-1 {
		p.tokenPosition++
	} else if p.tokenPosition == len(p.tokens)-1 {
		p.tokenPosition++
	} else {
		panic("Unexpected EOF")
	}
}

func (p *parser) consume(t tokens.TokenType) {
	if p.tokenPosition < len(p.tokens)-1 {
		if p.currentToken().Type != t {
			panic(fmt.Errorf("Unexpected token %q, expected %q", p.currentToken(), t))
		}
		p.nextToken()
	} else if p.tokenPosition == len(p.tokens)-1 {
		if p.currentToken().Type != t {
			panic(fmt.Errorf("Unexpected token %q, expected %q", p.currentToken(), t))
		}
		p.tokenPosition++
	} else {
		panic("Unexpected EOF")
	}
}

func (p *parser) expect(t tokens.TokenType) {
	if p.currentToken().Type != t {
		panic(fmt.Errorf("Unexpected token %q, expected %q", p.currentToken(), t))
	}
}

func (p *parser) currentToken() tokens.Token {
	if p.tokenPosition < len(p.tokens) {
		return p.tokens[p.tokenPosition]
	} else {
		panic("Unexpected EOF")
	}
}

func (p *parser) peekToken() tokens.Token {
	if p.tokenPosition < len(p.tokens)-1 {
		return p.tokens[p.tokenPosition + 1]
	} else {
		panic("Unexpected EOF")
	}
}

func (p *parser) parseStatement() nodes.Statement {
	statementParser, exists := p.statementParsers[p.currentToken().Type]
	if exists {
		return statementParser()
	} else {
		return nodes.ExpressionStatement {
			Expression: p.parseExpression(LOWEST),
		}
	}
}

func (p *parser) parseIf() nodes.Statement {
	n := &nodes.If {
		Location: p.currentToken().Location,
	}

	p.nextToken()

	n.Condition = p.parseExpression(LOWEST)

	n.Primary = p.parseStatement()

	if p.currentToken().Type == tokens.Else {
		p.nextToken()

		n.Alternative = p.parseStatement()
	}

	return n
}

func (p *parser) parseFor() nodes.Statement {
	n := &nodes.For {
		Location: p.currentToken().Location,
	}

	p.nextToken()

	n.PreStatement = p.parseStatement()

	p.consume(tokens.Semicolon)

	n.Condition = p.parseExpression(LOWEST)

	p.consume(tokens.Semicolon)

	n.PostStatement = p.parseStatement()

	n.Loop = p.parseStatement()

	return n
}

func (p *parser) parseAssignment() nodes.Statement {
	n := &nodes.Assignment {
		Location: p.currentToken().Location,
	}


	p.nextToken()


	n.Place = p.parseExpression(LOWEST)


	p.consume(tokens.Assignment)


	n.Value = p.parseExpression(LOWEST)


	return n
}

func (p *parser) parseScope() nodes.Statement {
	n := &nodes.Scope {
		Statements: make([]nodes.Statement, 0),
		Location: p.currentToken().Location,
	}

	p.nextToken()

	for p.currentToken().Type != tokens.CloseCurlyBracket {
		n.Statements = append(n.Statements, p.parseStatement())
	}

	p.nextToken()

	return n
}

func (p *parser) parseExpression(precedence Precedence) nodes.Expression {
	prefixParser, exists := p.prefixParsers[p.currentToken().Type]
	if !exists {
		// ERROR: can't parse prefix
		panic(fmt.Errorf("Failed to parse token %q, no prefix parser found", p.currentToken()))
	}


	left := prefixParser()

	for precedence < getPrecedence(p.currentToken().Type) {
		infixParser, exists := p.infixParsers[p.currentToken().Type]
		if !exists {
			return left
		}

		left = infixParser(left)
	}


	return left
}

func (p *parser) parseIntLiteral() nodes.Expression {
	defer p.nextToken()

	literal := strings.ToLower(p.currentToken().Literal)

	if strings.ContainsRune(literal, 'e') {
		vals := strings.SplitN(p.currentToken().Literal, "e", 2)
		coefficient, err := strconv.Atoi(vals[0])
		if err != nil {
			// ERROR: can't parse IntLiteral
			panic(fmt.Errorf("Failed to parse %q into IntLiteral, %s", p.currentToken(), err))
		}
		exponent, err := strconv.Atoi(vals[1])
		if err != nil {
			// ERROR: can't parse IntLiteral
			panic(fmt.Errorf("Failed to parse %q into IntLiteral, %s", p.currentToken(), err))
		}

		return &nodes.IntLiteral {
			Value: coefficient * int(math.Pow10(exponent)),
			Location: p.currentToken().Location,
		}
	} else {
		val, err := strconv.Atoi(p.currentToken().Literal)
		if err != nil {
			// ERROR: invalid IntLiteral token
			panic(fmt.Errorf("Failed to parse %q into IntLiteral, %s", p.currentToken(), err))
		}

		return &nodes.IntLiteral {
			Value: val,
			Location: p.currentToken().Location,
		}
	}
}

func (p *parser) parseFloatLiteral() nodes.Expression {
	defer p.nextToken()

	val, err := strconv.ParseFloat(p.currentToken().Literal, 32)
	if err != nil {
		// ERROR: can't parse FloatLiteral
			panic(fmt.Errorf("Failed to parse %q into FloatLiteral, %s", p.currentToken(), err))
	}

	return &nodes.FloatLiteral {
		Value: float32(val),
		Location: p.currentToken().Location,
	}
}

func (p *parser) parseStringLiteral() nodes.Expression {
	defer p.nextToken()

	return &nodes.StringLiteral {
		Value: strings.TrimPrefix(strings.TrimSuffix(p.currentToken().Literal, "\""), "\""),
		Location: p.currentToken().Location,
	}
}

func (p *parser) parseIdentifier() nodes.Expression {
	defer p.nextToken()

	// HACK
	if p.currentToken().Literal == "true" {
		return &nodes.BoolLiteral {
			Value: true,
			Location: p.currentToken().Location,
		}
	} else if p.currentToken().Literal == "false" {
		return &nodes.BoolLiteral {
			Value: false,
			Location: p.currentToken().Location,
		}
	} else {
		return &nodes.Identifier {
			Name: p.currentToken().Literal,
			Location: p.currentToken().Location,
		}
	}


}

func (p *parser) parsePrefixOperator() nodes.Expression {
	n := &nodes.UnaryOperator {
		Type: p.currentToken().Literal,
		Location: p.currentToken().Location,
	}

	p.nextToken()

	n.Operand = p.parseExpression(UNARY)

	return n
}

func (p *parser) parseOperator(left nodes.Expression) nodes.Expression {
	n := &nodes.Operator {
		Type: p.currentToken().Literal,
		Left: left,
		Location: p.currentToken().Location,
	}

	opPrecendence := getPrecedence(p.currentToken().Type)

	p.nextToken()

	n.Right = p.parseExpression(opPrecendence)

	return n
}

func (p *parser) parseCall(left nodes.Expression) nodes.Expression {
	n := &nodes.Call {
		Function: left,
		Arguments: make([]nodes.Expression, 0),
		Location: p.currentToken().Location,
	}

	p.nextToken()

	if p.currentToken().Type == tokens.CloseBracket {
		p.nextToken()
		return n
	}

	n.Arguments = append(n.Arguments, p.parseExpression(LOWEST))

	for ; p.currentToken().Type == tokens.Comma; p.nextToken() {
		n.Arguments = append(n.Arguments, p.parseExpression(LOWEST))
	}

	p.consume(tokens.CloseBracket)

	return n
}

func (p *parser) parseIndex(left nodes.Expression) nodes.Expression {
	n := &nodes.Index {
		Structure: left,
		Location: p.currentToken().Location,
	}

	p.nextToken()

	n.Index = p.parseExpression(LOWEST)

	p.consume(tokens.CloseSquareBracket)

	return n
}

func (p *parser) parseSubExpression() nodes.Expression {
	p.nextToken()

	n := p.parseExpression(LOWEST)

	p.consume(tokens.CloseBracket)

	return n
}
