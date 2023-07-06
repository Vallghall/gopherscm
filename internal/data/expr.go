package data

import (
	"encoding/json"
	"errors"
)

// Expr - expressions supported by the program
type Expr uint

// ErrUnsupportedExprKind - error unsupported expression kind
var ErrUnsupportedExprKind = errors.New("unsupported expression kind")

// Expr enum
const (
	// Literal - expression which is evaluated into itself
	Literal Expr = iota
	// CallExpr - function with a list of arguments
	CallExpr
	// VariableRef - variable which value should be known
	// before evaluation during interpretation
	VariableRef
	// DefineExpr - expression that creates a new lexical scope
	// and binds variable with an expression
	DefineExpr
	// Function - node containing function body
	Function
	// Root - AST root unique expressions kind
	Root = 9999
)

func (e Expr) MarshalJSON() ([]byte, error) {
	switch e {
	case Literal:
		return json.Marshal("Literal")
	case CallExpr:
		return json.Marshal("CallExpr")
	case DefineExpr:
		return json.Marshal("DefineExpr")
	case Root:
		return json.Marshal("Root")
	default:
		return nil, ErrUnsupportedExprKind
	}
}
