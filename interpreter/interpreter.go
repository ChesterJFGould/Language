package main

import (
	"fmt"
	"strconv"

	"../nodes"
)

type interpreter struct {
	symbolTable []map[string] nodes.Expression
}

func newInterpreter() *interpreter {
	return &interpreter {
		symbolTable: make([]map[string] nodes.Expression, 0),
	}
}

func (i *interpreter) retrieveSymbol(s string) (nodes.Expression, bool) {
	for _, m := range i.symbolTable {
		e, exists := m[s]
		if exists {
			return e, true
		}
	}

	return nil, false
}

func (i *interpreter) storeSymbol(name string, e nodes.Expression) {
	for n := range i.symbolTable {
		_, exists := i.symbolTable[n][name]
		if exists {
			i.symbolTable[n][name] = e
			return
		}
	}
	if len(i.symbolTable) > 0 {
		i.symbolTable[len(i.symbolTable)-1][name] = e
	} else {
		panic(fmt.Sprintf("Attempted to store symbol without a scope at %s", e.GetLocation()))
	}
}

func (i *interpreter) newScope() {
	i.symbolTable = append(i.symbolTable, make(map[string] nodes.Expression))
}

func (i *interpreter) deleteScope() {
	if len(i.symbolTable) > 0 {
		i.symbolTable = i.symbolTable[:len(i.symbolTable)-1]
	}
}

func (i *interpreter) interpretExpression(e nodes.Expression) nodes.Expression {
	switch e := e.(type) {
	case *nodes.IntLiteral:
		return e
	case *nodes.FloatLiteral:
		return e
	case *nodes.StringLiteral:
		return e
	case *nodes.Identifier:
		val, exists := i.retrieveSymbol(e.Name)
		if exists {
			return val
		} else {
			panic(fmt.Sprintf("Undeclared variable %q at %s", e.Name, e.Location))
		}
	case *nodes.Call:
		switch function := e.Function.(type) {
		case *nodes.Identifier:
			switch function.Name {
			case "println":
				if len(e.Arguments) > 1 {
					panic(fmt.Sprintf("Too many arguments in call to println at %s", e.Location))
				} else if len(e.Arguments) == 0 {
					panic(fmt.Sprintf("Too few arguments in call to println at %s", e.Location))
				} else {
					switch s := i.interpretExpression(e.Arguments[0]); s := s.(type) {
					case *nodes.StringLiteral:
						fmt.Println(s.Value)
						return &nodes.Void {
							Location: e.Location,
						}
					default:
						panic(fmt.Sprintf("Cannot use token at %s as type string in call to println", s.GetLocation()))
					}
				}
			case "string":
				if len(e.Arguments) > 1 {
					panic(fmt.Sprintf("Too many arguments in call to string at %s", e.Location))
				} else if len(e.Arguments) == 0 {
					panic(fmt.Sprintf("Too few arguments in call to string at %s", e.Location))
				} else {
					switch v := i.interpretExpression(e.Arguments[0]); v := v.(type) {
					case *nodes.StringLiteral:
						return e
					case *nodes.IntLiteral:
						return &nodes.StringLiteral {
							Value: strconv.Itoa(v.Value),
							Location: e.Location,
						}
					case *nodes.FloatLiteral:
						return &nodes.StringLiteral {
							Value: strconv.FormatFloat(float64(v.Value), 'E', -1, 32),
							Location: e.Location,
						}
					default:
						panic(fmt.Sprintf("Invalid argument in call to string at %s", v.GetLocation()))
					}
				}
			default:
				panic("Sorry Chief, still working on non built-in functions")
			}
		default:
			panic(fmt.Sprintf("Cannot use token at %s as function in Call", e.Location))
		}
	case *nodes.Index:
		panic("Sorry Chief, still working on supporting arrays")
	case *nodes.Operator:
		left := i.interpretExpression(e.Left)
		right := i.interpretExpression(e.Right)

		switch left := left.(type) {
		case *nodes.IntLiteral:
			switch right := right.(type) {
			case *nodes.IntLiteral:
				switch e.Type {
				case "+":
					return &nodes.IntLiteral {
						Value: left.Value + right.Value,
						Location: e.Location,
					}
				case "-":
					return &nodes.IntLiteral {
						Value: left.Value - right.Value,
						Location: e.Location,
					}
				case "*":
					return &nodes.IntLiteral {
						Value: left.Value * right.Value,
						Location: e.Location,
					}
				case "/":
					return &nodes.IntLiteral {
						Value: left.Value / right.Value,
						Location: e.Location,
					}
				case "<":
					return &nodes.BoolLiteral {
						Value: left.Value < right.Value,
						Location: e.Location,
					}
				case ">":
					return &nodes.BoolLiteral {
						Value: left.Value > right.Value,
						Location: e.Location,
					}
				case "==":
					return &nodes.BoolLiteral {
						Value: left.Value == right.Value,
						Location: e.Location,
					}
				default:
					panic(fmt.Sprintf("Operator %q at %s not defined on Float", e.Type, e.Location))
				}
			default:
				panic(fmt.Sprintf("Mismatched types on Operator at %s", e.Location))
			}
		case *nodes.FloatLiteral:
			switch right := right.(type) {
			case *nodes.FloatLiteral:
				switch e.Type {
				case "+":
					return &nodes.FloatLiteral {
						Value: left.Value + right.Value,
						Location: e.Location,
					}
				case "-":
					return &nodes.FloatLiteral {
						Value: left.Value - right.Value,
						Location: e.Location,
					}
				case "*":
					return &nodes.FloatLiteral {
						Value: left.Value * right.Value,
						Location: e.Location,
					}
				case "/":
					return &nodes.FloatLiteral {
						Value: left.Value / right.Value,
						Location: e.Location,
					}
				case "<":
					return &nodes.BoolLiteral {
						Value: left.Value < right.Value,
						Location: e.Location,
					}
				case ">":
					return &nodes.BoolLiteral {
						Value: left.Value > right.Value,
						Location: e.Location,
					}
				case "==":
					return &nodes.BoolLiteral {
						Value: left.Value == right.Value,
						Location: e.Location,
					}
				default:
					panic(fmt.Sprintf("Operator %q at %s not defined on Float", e.Type, e.Location))
				}
			default:
				panic(fmt.Sprintf("Mismatched types on Operator at %s", e.Location))
			}
		case *nodes.StringLiteral:
			switch right := right.(type) {
			case *nodes.StringLiteral:
				switch e.Type {
				case "+":
					return &nodes.StringLiteral {
						Value: left.Value + right.Value,
						Location: e.Location,
					}
				case "==":
					return &nodes.BoolLiteral {
						Value: left.Value == right.Value,
						Location: e.Location,
					}
				default:
					panic(fmt.Sprintf("Operator %q at %s not defined on String", e.Type, e.Location))
				}
			default:
				panic(fmt.Sprintf("Mismatched types on Operator at %s", e.Location))
			}
		case *nodes.BoolLiteral:
			switch right := right.(type) {
			case *nodes.BoolLiteral:
				switch e.Type {
				case "==":
					return &nodes.BoolLiteral {
						Value: left.Value == right.Value,
						Location: e.Location,
					}
				default:
					panic(fmt.Sprintf("Operator %q at %s not defined on Bool", e.Type, e.Location))
				}
			default:
				panic(fmt.Sprintf("Mismatched types on Operator at %s", e.Location))
			}
		}
	case *nodes.UnaryOperator:
		operand := i.interpretExpression(e.Operand)

		switch operand := operand.(type) {
		case *nodes.IntLiteral:
			switch e.Type {
			case "-":
				return &nodes.IntLiteral {
					Value: -operand.Value,
					Location: e.Location,
				}
			default:
				panic(fmt.Sprintf("Operator %q at %s not defined on Int", e.Type, e.Location))
			}
		case *nodes.FloatLiteral:
			switch e.Type {
			case "-":
				return &nodes.FloatLiteral {
					Value: -operand.Value,
					Location: e.Location,
				}
			default:
				panic(fmt.Sprintf("Operator %q at %s not defined on Float", e.Type, e.Location))
			}
		case *nodes.StringLiteral:
			switch e.Type {
			default:
				panic(fmt.Sprintf("Operator %q at %s not defined on String", e.Type, e.Location))
			}
		case *nodes.BoolLiteral:
			switch e.Type {
			default:
				panic(fmt.Sprintf("Operator %q at %s not defined on Bool", e.Type, e.Location))
			}
		}
	}

	panic(fmt.Sprintf("Failed to match node at %s", e.GetLocation()))
}


func (i *interpreter) interpretStatement(s nodes.Statement) {
	switch s := s.(type) {
	case *nodes.If:
		i.newScope()
		condition := i.interpretExpression(s.Condition)
		switch condition := condition.(type) {
		case *nodes.BoolLiteral:
			if condition.Value {
				i.interpretStatement(s.Primary)
			} else if s.Alternative != nil {
				i.interpretStatement(s.Alternative)
			}
		default:
			panic(fmt.Sprintf("Non bool expression used as condition at %s", s.Location))
		}
		i.deleteScope()
	case *nodes.For:
		i.newScope()
		i.interpretStatement(s.PreStatement)

		for {
			var done bool
			condition := i.interpretExpression(s.Condition)
			switch condition := condition.(type) {
			case *nodes.BoolLiteral:
				done = !condition.Value
			default:
				panic(fmt.Sprintf("Non bool expression used as condition at %s", s.Location))
			}
			if done {
				break
			}
			i.interpretStatement(s.Loop)
			i.interpretStatement(s.PostStatement)
		}
		i.deleteScope()
	case *nodes.Assignment:
		switch place := s.Place.(type) {
		case *nodes.Identifier:
			i.storeSymbol(place.Name, i.interpretExpression(s.Value))
		default:
			panic(fmt.Sprintf("Invalid left hand side of assignment at %s", s.Location))
		}
	case *nodes.Scope:
		i.newScope()
		for _, statement := range s.Statements {
			i.interpretStatement(statement)
		}
		i.deleteScope()
	case *nodes.ExpressionStatement:
		//println("ExpressionStatement")
		i.interpretExpression(s.Expression)
	}
}
