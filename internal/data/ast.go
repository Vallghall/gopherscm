package data

import (
	"github.com/Vallghall/gopherscm/internal/core"
)

// AST - represents Scheme program structure
type AST struct {
	Ctx      *Context
	Token    *Token
	Kind     Expr
	Subtrees []*AST
}

// ASTRoot - constructor for AST
func ASTRoot() *AST {
	return &AST{
		Kind:     Root,
		Ctx:      NewContext(core.DefaultDefinitions()),
		Subtrees: make([]*AST, 0),
	}
}

// Nest - AST node constructor
func (ast *AST) Nest(t *Token) *AST {
	node := &AST{
		Token:    t,
		Ctx:      ast.Ctx,
		Subtrees: make([]*AST, 0),
	}

	if t.Value() == "define" {
		node.Kind = DefineExpr
	} else {
		node.Kind = CallExpr // all nested forms have functions as the first elem
	}

	ast.Subtrees = append(ast.Subtrees, node)

	return node
}

// Identifier - returns value of stored Token
func (ast *AST) Identifier() string {
	return ast.Token.Value()
}

// Add - AST node constructor
func (ast *AST) Add(t *Token) *AST {
	var e Expr
	switch t.Type() {
	case Int, Float, String:
		e = Literal
	case Id:
		e = VariableRef
	default: // fill in later
	}

	node := &AST{
		Ctx:      ast.Ctx,
		Token:    t,
		Kind:     e,
		Subtrees: make([]*AST, 0),
	}

	ast.Subtrees = append(ast.Subtrees, node)

	return node
}
