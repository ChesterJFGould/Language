package nodes

import (
	"fmt"
	"strings"
	"bufio"
	"strconv"

	"../location"
)

////////////////
// Statements //
////////////////

type Statement interface {
	statementNode()
	PrintTree(indent string, last bool)
	String() string
	GetLocation() location.Location
}

var statementScannerParsers map[string] func (s *bufio.Scanner) (Statement, error)

func init() {
	statementScannerParsers = make(map[string]func(s *bufio.Scanner)(Statement, error))

	statementScannerParsers["If"] = IfFromScanner
	statementScannerParsers["For"] = ForFromScanner
	statementScannerParsers["Assignment"] = AssignmentFromScanner
	statementScannerParsers["Scope"] = ScopeFromScanner
}

func StatementFromScanner(s *bufio.Scanner) (Statement, error) {
	vals := strings.Split(s.Text(), " ")

	parser, exists := statementScannerParsers[vals[0]]
	if exists {
		return parser(s)
	} else {
		return ExpressionStatementFromScanner(s)
	}
}

type If struct {
	Condition Expression
	Primary Statement
	Alternative Statement
	location.Location
}

func (i If) statementNode() {}

func (i If) PrintTree(indent string, last bool) {
	fmt.Print(indent)

	if last {
		fmt.Print("\\-")
		indent += "  "
	} else {
		fmt.Print("|-")
		indent += "| "
	}

	fmt.Print("If\n")

	if i.Alternative == nil {
		i.Primary.PrintTree(indent, true)
	} else {
		i.Primary.PrintTree(indent, false)
		i.Alternative.PrintTree(indent, true)
	}
}

func (i If) String() string {
	var b strings.Builder

	b.WriteString("If if ")

	if i.Alternative == nil {
		b.WriteString("2 "+i.Location.String()+"\n")
		b.WriteString(i.Condition.String()+"\n")
		b.WriteString(i.Primary.String())
	} else {
		b.WriteString("3 "+i.Location.String()+"\n")
		b.WriteString(i.Condition.String()+"\n")
		b.WriteString(i.Primary.String()+"\n")
		b.WriteString(i.Alternative.String())
	}

	return b.String()
}

func (i If) GetLocation() location.Location {
	return i.Location
}

func IfFromScanner(s *bufio.Scanner) (Statement, error) {
	vals := strings.Split(s.Text(), " ")
	if vals[0] != "If" {
		return nil, fmt.Errorf("Failed to parse %q into If", s.Text())
	}

	loc, err := location.LocationFromString(strings.Join(vals[3:], " "))
	if err != nil {
		return nil, err
	}

	i := &If {
		Location: loc,
	}

	ok := s.Scan()
	if !ok {
		return nil, fmt.Errorf("Failed to parse If from scanner: EOF") 
	}

	i.Condition, err = ExpressionFromScanner(s)
	if err != nil {
		return nil, err
	}

	ok = s.Scan()
	if !ok {
		return nil, fmt.Errorf("Failed to parse If from scanner: EOF")
	}

	i.Primary, err = StatementFromScanner(s)
	if err != nil {
		return nil, err
	}

	if vals[3] == "2" {
		return i, nil
	}

	ok = s.Scan()
	if !ok {
		return nil, fmt.Errorf("Failed to parse If from scanner: EOF")
	}

	i.Alternative, err = StatementFromScanner(s)
	if err != nil  {
		return nil, err
	}

	return i, nil
}


type For struct {
	PreStatement Statement
	Condition Expression
	PostStatement Statement
	Loop Statement
	location.Location
}

func (f For) statementNode() {}

func (f For) PrintTree(indent string, last bool) {
	fmt.Print(indent)

	if last {
		fmt.Print("\\-")
		indent += "  "
	} else {
		fmt.Print("|-")
		indent += "| "
	}

	fmt.Print("For\n")

	f.PreStatement.PrintTree(indent, false)
	f.Condition.PrintTree(indent, false)
	f.PostStatement.PrintTree(indent, false)
	f.Loop.PrintTree(indent, true)
}

func (f For) String() string {
	var b strings.Builder

	b.WriteString("For for 4 "+f.Location.String()+"\n")

	b.WriteString(f.PreStatement.String()+"\n")
	b.WriteString(f.Condition.String()+"\n")
	b.WriteString(f.PostStatement.String()+"\n")
	b.WriteString(f.Loop.String())

	return b.String()
}

func (f For) GetLocation() location.Location {
	return f.Location
}

func ForFromScanner(s *bufio.Scanner) (Statement, error) {
	vals := strings.Split(s.Text(), " ")

	if vals[0] != "For" {
		return nil, fmt.Errorf("Failed to parse %q into For", s.Text())
	}

	loc, err := location.LocationFromString(strings.Join(vals[3:], " "))
	if err != nil {
		return nil, err
	}

	f := &For {
		Location: loc,
	}

	ok := s.Scan()
	if !ok {
		return nil, fmt.Errorf("Failed to parse For from scanner: EOF")
	}

	f.PreStatement, err = StatementFromScanner(s)
	if err != nil {
		return nil, err
	}

	ok = s.Scan()
	if !ok {
		return nil, fmt.Errorf("Failed to parse For from scanner: EOF")
	}

	f.Condition, err = ExpressionFromScanner(s)
	if err != nil {
		return nil, err
	}

	ok = s.Scan()
	if !ok {
		return nil, fmt.Errorf("Failed to parse For from scanner: EOF")
	}

	f.PostStatement, err = StatementFromScanner(s)
	if err != nil {
		return nil, err
	}

	ok = s.Scan()
	if !ok {
		return nil, fmt.Errorf("Failed to parse For from scanner: EOF")
	}

	f.Loop, err = StatementFromScanner(s)
	if err != nil {
		return nil, err
	}

	return f, nil
}

type Assignment struct {
	Place Expression
	Value Expression
	location.Location
}

func (a Assignment) statementNode() {}

func (a Assignment) PrintTree(indent string, last bool) {
	fmt.Print(indent)

	if last {
		fmt.Print("\\-")
		indent += "  "
	} else {
		fmt.Print("|-")
		indent += "| "
	}

	fmt.Print("=\n")

	a.Place.PrintTree(indent, false)
	a.Value.PrintTree(indent, true)
}

func (a Assignment) String() string {
	var b strings.Builder

	b.WriteString("Assignment = 2 "+a.Location.String()+"\n")
	b.WriteString(a.Place.String()+"\n")
	b.WriteString(a.Value.String())

	return b.String()
}

func (a Assignment) GetLocation() location.Location {
	return a.Location
}

func AssignmentFromScanner(s *bufio.Scanner) (Statement, error) {
	vals := strings.Split(s.Text(), " ")

	if vals[0] != "Assignment" {
		return nil, fmt.Errorf("Failed to parse %q into Assignment", s.Text())
	}

	loc, err := location.LocationFromString(strings.Join(vals[3:], " "))
	if err != nil {
		return nil, err
	}

	a := &Assignment {
		Location: loc,
	}

	ok := s.Scan()
	if !ok {
		return nil, fmt.Errorf("Failed to parse Assignment from scanner: EOF")
	}

	a.Place, err = ExpressionFromScanner(s)
	if err != nil {
		return nil, err
	}

	ok = s.Scan()
	if !ok {
		return nil, fmt.Errorf("Failed to parse Assignment from scanner: EOF")
	}

	a.Value, err = ExpressionFromScanner(s)
	if err != nil {
		return nil, err
	}

	return a, nil
}

type Scope struct {
	Statements []Statement
	location.Location
}

func (s Scope) statementNode() {}

func (s Scope) PrintTree(indent string, last bool) {
	fmt.Print(indent)
	if last {
		fmt.Print("\\-")
		indent += "  "
	} else {
		fmt.Print("|-")
		indent += "| "
	}

	fmt.Print("Scope\n")

	for _, s := range s.Statements[:len(s.Statements)-1] {
		s.PrintTree(indent, false)
	}

	s.Statements[len(s.Statements)-1].PrintTree(indent, true)
}

func (s Scope) String() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("Scope { %d %s\n", len(s.Statements), s.Location))
	for _, s := range s.Statements[:len(s.Statements)-1] {
		b.WriteString(s.String()+"\n")
	}
	b.WriteString(s.Statements[len(s.Statements)-1].String())

	return b.String()
}

func (s Scope) GetLocation() location.Location {
	return s.Location
}

func ScopeFromScanner(s *bufio.Scanner) (Statement, error) {
	vals := strings.Split(s.Text(), " ")

	if vals[0] != "Scope" {
		return nil, fmt.Errorf("Failed to parse %q into Scope", s.Text())
	}

	loc, err := location.LocationFromString(strings.Join(vals[3:], " "))
	if err != nil {
		return nil, err
	}

	scope := &Scope {
		Statements: make([]Statement, 0),
		Location: loc,
	}

	numSubnodes, err := strconv.Atoi(vals[2])
	if err != nil {
		return nil, fmt.Errorf("Failed to parse Scope from scanner: %s", err)
	}

	for i := 0; i < numSubnodes; i++ {
		ok := s.Scan()
		if !ok {
			return nil, fmt.Errorf("Failed to parse Scope from scanner: EOF")
		}

		subStatement, err := StatementFromScanner(s)
		if err != nil {
			return nil, err
		}

		scope.Statements = append(scope.Statements, subStatement)
	}

	return scope, nil
}


type ExpressionStatement struct {
	Expression
}

func (e ExpressionStatement) statementNode() {}

func (e ExpressionStatement) PrintTree(indent string, last bool) {
	e.Expression.PrintTree(indent, last)
}

func (e ExpressionStatement) String() string {
	return e.Expression.String()
}

func (e ExpressionStatement) GetLocation() location.Location {
	return e.Expression.GetLocation()
}

func ExpressionStatementFromScanner(s *bufio.Scanner) (Statement, error) {
	e, err := ExpressionFromScanner(s)
	if err != nil {
		return nil, err
	}

	return &ExpressionStatement{e}, nil
}

/////////////////
// Expressions //
/////////////////

type Expression interface {
	expressionNode()
	PrintTree(indent string, last bool)
	String() string
	GetLocation() location.Location
}

var expressionScannerParsers map[string] func(s *bufio.Scanner) (Expression, error)

func init() {
	expressionScannerParsers = make(map[string]func(s *bufio.Scanner)(Expression, error))

	expressionScannerParsers["IntLiteral"] = IntLiteralFromScanner
	expressionScannerParsers["FloatLiteral"] = FloatLiteralFromScanner
	expressionScannerParsers["StringLiteral"] = StringLiteralFromScanner
	expressionScannerParsers["Identifier"] = IdentifierFromScanner
	expressionScannerParsers["Call"] = CallFromScanner
	expressionScannerParsers["Index"] = IndexFromScanner
	expressionScannerParsers["Operator"] = OperatorFromScanner
	expressionScannerParsers["UnaryOperator"] = UnaryOperatorFromScanner
}

func ExpressionFromScanner(s *bufio.Scanner) (Expression, error) {
	vals := strings.Split(s.Text(), " ")

	parser, exists := expressionScannerParsers[vals[0]]
	if exists {
		return parser(s)
	} else {
		return nil, fmt.Errorf("Failed to parse %q into Expression", s.Text())
	}
}

type IntLiteral struct {
	Value int
	location.Location
}

func (i IntLiteral) expressionNode() {}

func (i IntLiteral) PrintTree(indent string, last bool) {
	fmt.Print(indent)

	if last {
		fmt.Print("\\-")
	} else {
		fmt.Print("|-")
	}

	fmt.Printf("%d\n", i.Value)
}

func (i IntLiteral) String() string {
	return "IntLiteral "+fmt.Sprintf("%d", i.Value)+" 0 "+i.Location.String()
}

func (i IntLiteral) GetLocation() location.Location {
	return i.Location
}

func IntLiteralFromScanner(s *bufio.Scanner) (Expression, error) {
	vals := strings.Split(s.Text(), " ")

	if vals[0] != "IntLiteral" {
		return nil, fmt.Errorf("Failed to parse %q into IntLiteral", s.Text())
	}

	loc, err := location.LocationFromString(strings.Join(vals[3:], " "))
	if err != nil {
		return nil, fmt.Errorf("Failed to parse IntLiteral from scanner: %s", err)
	}

	i := &IntLiteral {
		Location: loc,
	}

	i.Value, err = strconv.Atoi(vals[1])
	if err != nil {
		return nil, err
	}

	return i, nil
}

type FloatLiteral struct {
	Value float32
	location.Location
}

func (f FloatLiteral) expressionNode() {}

func (f FloatLiteral) PrintTree(indent string, last bool) {
	fmt.Print(indent)

	if last {
		fmt.Print("\\-")
	} else {
		fmt.Print("|-")
	}

	fmt.Printf("%g\n", f.Value)
}

func (f FloatLiteral) String() string {
	return "FloatLiteral "+fmt.Sprintf("%g", f.Value)+" 0 "+f.Location.String()
}

func (f FloatLiteral) GetLocation() location.Location {
	return f.Location
}

func FloatLiteralFromScanner(s *bufio.Scanner) (Expression, error) {
	vals := strings.Split(s.Text(), " ")

	if vals[0] != "FloatLiteral" {
		return nil, fmt.Errorf("Failed to parse %q into FloatLiteral", s.Text())
	}

	loc, err := location.LocationFromString(strings.Join(vals[3:], " "))
	if err != nil {
		return nil, fmt.Errorf("Failed to parse FloatLiteral from scanner: %s", err)
	}

	f := &FloatLiteral {
		Location: loc,
	}

	val, err := strconv.ParseFloat(vals[1], 32)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse FloatLiteral from scanner: %s", err)
	}

	f.Value = float32(val)

	return f, nil
}

type StringLiteral struct {
	Value string
	location.Location
}

func (s StringLiteral) expressionNode() {}

func (s StringLiteral) PrintTree(indent string, last bool) {
	fmt.Print(indent)

	if last {
		fmt.Print("\\-")
	} else {
		fmt.Print("|-")
	}

	fmt.Printf("%q\n", s.Value)
}

func (s StringLiteral) String() string {
	return "StringLiteral "+s.Value+" 0 "+s.Location.String()
}

func (s StringLiteral) GetLocation() location.Location {
	return s.Location
}

func StringLiteralFromScanner(s *bufio.Scanner) (Expression, error) {
	vals := strings.Split(s.Text(), " ")

	if vals[0] != "StringLiteral" {
		return nil, fmt.Errorf("Failed to parse %q into StringLiteral", s.Text())
	}

	loc, err := location.LocationFromString(strings.Join(vals[len(vals)-3:], " "))
	if err != nil {
		return nil, fmt.Errorf("Failed to parse StringLiteral from scanner: %s", err)
	}

	return &StringLiteral {
		Value: strings.Join(vals[1:len(vals)-4], " "),
		Location: loc,
	}, nil
}

type BoolLiteral struct {
	Value bool
	location.Location
}

func (b BoolLiteral) expressionNode() {}

func (b BoolLiteral) PrintTree(indent string, last bool) {
	fmt.Print(indent)

	if last {
		fmt.Print("\\-")
	} else {
		fmt.Print("|-")
	}

	fmt.Printf("%t\n", b.Value)
}

func (b BoolLiteral) String() string {
	return "BoolLiteral "+fmt.Sprintf("%t", b.Value)+" 0 "+b.Location.String()
}

func (b BoolLiteral) GetLocation() location.Location {
	return b.Location
}

func BoolLiteralFromScanner(s *bufio.Scanner) (Expression, error) {
	vals := strings.Split(s.Text(), " ")

	if vals[0] != "BoolLiteral" {
		return nil, fmt.Errorf("Failed to parse %q into BoolLiteral", s.Text())
	}

	loc, err := location.LocationFromString(strings.Join(vals[4:], " "))
	if err != nil {
		return nil, err
	}

	switch vals[1] {
	case "true":
		return &BoolLiteral {
			Value: true,
			Location: loc,
		}, nil
	case "false":
		return &BoolLiteral {
			Value: false,
			Location: loc,
		}, nil
	default:
		return nil, fmt.Errorf("Failed to parse %q into BoolLiteral", s.Text())
	}
}

type Identifier struct {
	Name string
	location.Location
}

func (i Identifier) expressionNode() {}

func (i Identifier) PrintTree(indent string, last bool) {
	fmt.Print(indent)

	if last {
		fmt.Print("\\-")
	} else {
		fmt.Print("|-")
	}

	fmt.Printf("%s\n", i.Name)
}

func (i Identifier) String() string {
	return "Identifier "+i.Name+" 0 "+i.Location.String()
}

func (i Identifier) GetLocation() location.Location {
	return i.Location
}

func IdentifierFromScanner(s *bufio.Scanner) (Expression, error) {
	vals := strings.Split(s.Text(), " ")

	if vals[0] != "Identifier" {
		return nil, fmt.Errorf("Failed to parse %q into Identifier", s.Text())
	}

	loc, err := location.LocationFromString(strings.Join(vals[3:], " "))
	if err != nil {
		return nil, fmt.Errorf("Failed to parse Identifier from scanner: %s", err)
	}

	return &Identifier {
		Name: vals[1],
		Location: loc,
	}, nil
}

type Call struct {
	Function Expression
	Arguments []Expression
	location.Location
}

func (c Call) expressionNode() {}

func (c Call) PrintTree(indent string, last bool) {
	fmt.Print(indent)

	if last {
		fmt.Print("\\-")
		indent += "  "
	} else {
		fmt.Print("|-")
		indent += "| "
	}

	fmt.Print("Call\n")

	c.Function.PrintTree(indent, false)

	for _, a := range c.Arguments[:len(c.Arguments)-1] {
		a.PrintTree(indent, false)
	}

	c.Arguments[len(c.Arguments)-1].PrintTree(indent, true)
}

func (c Call) String() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("Call ( %d %s\n", 1+len(c.Arguments), c.Location))
	b.WriteString(c.Function.String()+"\n")
	for _, a := range c.Arguments[:len(c.Arguments)-1] {
		b.WriteString(a.String()+"\n")
	}
	b.WriteString(c.Arguments[len(c.Arguments)-1].String())

	return b.String()
}

func (c Call) GetLocation() location.Location {
	return c.Location
}

func CallFromScanner(s *bufio.Scanner) (Expression, error) {
	vals := strings.Split(s.Text(), " ")

	if vals[0] != "Call" {
		return nil, fmt.Errorf("Failed to parse %q into Call", s.Text())
	}

	loc, err := location.LocationFromString(strings.Join(vals[len(vals)-3:], " "))
	if err != nil {
		return nil, fmt.Errorf("Failed to parse Call from scanner: %s", err)
	}

	c := &Call {
		Arguments: make([]Expression, 0),
		Location: loc,
	}

	ok := s.Scan()
	if !ok {
		return nil, fmt.Errorf("Failed to parse Call from scanner: EOF")
	}

	c.Function, err = ExpressionFromScanner(s)
	if err != nil {
		return nil, err
	}

	numSubnodes, err := strconv.Atoi(vals[2])
	if err != nil {
		return nil, fmt.Errorf("Failed to parse Call from scanner: %s", err)
	}

	for i := 0; i < numSubnodes - 1; i++ {
		ok = s.Scan()
		if !ok {
			return nil, fmt.Errorf("Failed to parse Call from scanner: EOF")
		}

		arg, err := ExpressionFromScanner(s)
		if err != nil {
			return nil, err
		}

		c.Arguments = append(c.Arguments, arg)
	}

	return c, nil
}

type Index struct {
	Structure Expression
	Index Expression
	location.Location
}

func (i Index) expressionNode() {}

func (i Index) PrintTree(indent string, last bool) {
	fmt.Print(indent)

	if last {
		fmt.Print("\\-")
		indent += "  "
	} else {
		fmt.Print("|-")
		indent += "| "
	}

	fmt.Print("Index\n")

	i.Structure.PrintTree(indent, false)

	i.Index.PrintTree(indent, true)
}

func (i Index) String() string {
	var b strings.Builder

	b.WriteString("Index [ 2 "+i.Location.String()+"\n")
	b.WriteString(i.Structure.String()+"\n")
	b.WriteString(i.Index.String())

	return b.String()
}

func (i Index) GetLocation() location.Location {
	return i.Location
}

func IndexFromScanner(s *bufio.Scanner) (Expression, error) {
	vals := strings.Split(s.Text(), " ")

	if vals[0] != "Index" {
		return nil, fmt.Errorf("Failed to parse %q into Call", s.Text())
	}

	loc, err := location.LocationFromString(strings.Join(vals[3:], " "))
	if err != nil {
		return nil, fmt.Errorf("Failed to parse Call from scanner: %s", err)
	}

	i := &Index {
		Location: loc,
	}

	ok := s.Scan()
	if !ok {
		return nil, fmt.Errorf("Failed to parse Call from scanner: EOF")
	}

	i.Structure, err = ExpressionFromScanner(s)
	if err != nil {
		return nil, err
	}

	ok = s.Scan()
	if !ok {
		return nil, fmt.Errorf("Failed to parse Call from scanner: EOF")
	}

	i.Index, err = ExpressionFromScanner(s)
	if err != nil {
		return nil, err
	}

	return i, nil
}

type Operator struct {
	Type string
	Left Expression
	Right Expression
	location.Location
}

func (o Operator) expressionNode() {}

func (o Operator) PrintTree(indent string, last bool) {
	fmt.Print(indent)

	if last {
		fmt.Print("\\-")
		indent += "  "
	} else {
		fmt.Print("|-")
		indent += "| "
	}

	fmt.Printf("%s\n", o.Type)

	o.Left.PrintTree(indent, false)

	o.Right.PrintTree(indent, true)
}

func (o Operator) String() string {
	var b strings.Builder

	b.WriteString("Operator "+o.Type+" 2 "+o.Location.String()+"\n")
	b.WriteString(o.Left.String()+"\n")
	b.WriteString(o.Right.String())

	return b.String()
}

func (o Operator) GetLocation() location.Location {
	return o.Location
}

func OperatorFromScanner(s *bufio.Scanner) (Expression, error) {
	vals := strings.Split(s.Text(), " ")

	if vals[0] != "Operator" {
		return nil, fmt.Errorf("Failed to parse %q into Operator", s.Text())
	}

	loc, err := location.LocationFromString(strings.Join(vals[3:], " "))
	if err != nil {
		return nil, fmt.Errorf("Failed to parse Operator from scanner: %s", err)
	}

	o := &Operator {
		Type: vals[1],
		Location: loc,
	}

	ok := s.Scan()
	if !ok {
		return nil, fmt.Errorf("Failed to parse Operator from scanner: EOF")
	}

	o.Left, err = ExpressionFromScanner(s)
	if err != nil {
		return nil, err
	}

	ok = s.Scan()
	if !ok {
		return nil, fmt.Errorf("Failed to parse Operator from scanner: EOF")
	}

	o.Right, err = ExpressionFromScanner(s)
	if err != nil {
		return nil, err
	}

	return o, nil
}

type UnaryOperator struct {
	Type string
	Operand Expression
	location.Location
}

func (u UnaryOperator) expressionNode() {}

func (u UnaryOperator) PrintTree(indent string, last bool) {
	fmt.Print(indent)

	if last {
		fmt.Print("\\-")
		indent += "  "
	} else {
		fmt.Print("|-")
		indent += "| "
	}

	fmt.Printf("%s\n", u.Type)

	u.Operand.PrintTree(indent, true)
}

func (u UnaryOperator) String() string {
	var b strings.Builder

	b.WriteString("UnaryOperator "+u.Type+" 1 "+u.Location.String()+"\n")
	b.WriteString(u.Operand.String())

	return b.String()
}

func (u UnaryOperator) GetLocation() location.Location {
	return u.Location
}

func UnaryOperatorFromScanner(s *bufio.Scanner) (Expression, error) {
	vals := strings.Split(s.Text(), " ")

	if vals[0] != "UnaryOperator" {
		return nil, fmt.Errorf("Failed to parse %q into UnaryOperator", s.Text())
	}

	loc, err := location.LocationFromString(strings.Join(vals[3:], " "))
	if err != nil {
		return nil, fmt.Errorf("Failed to parse UnaryOperator from scanner: %s", err)
	}

	u := &UnaryOperator {
		Type: vals[1],
		Location: loc,
	}

	ok := s.Scan()
	if !ok {
		return nil, fmt.Errorf("Failed to parse UnaryOperator from scanner: EOF")
	}

	u.Operand, err = ExpressionFromScanner(s)
	if err != nil {
		return nil, err
	}

	return u, nil
}

type Void struct {
	location.Location
}

func (v Void) expressionNode() {}

func (v Void) PrintTree(indent string, last bool) {}

func (v Void) GetLocation() location.Location {
	return v.Location
}

func (v Void) String() string {
	return ""
}
