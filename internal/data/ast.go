package data

import (
	"errors"
	"fmt"
	"github.com/Vallghall/gopherscm/internal/core"
	"github.com/Vallghall/gopherscm/internal/core/types"
	"strconv"
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
		Ctx:      ast.Ctx.Spawn(),
		Token:    t,
		Kind:     CallExpr, // all nested forms have functions as the first elem
		Subtrees: make([]*AST, 0),
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

// Eval - evaluate AST to a single value
func (ast *AST) Eval() (res types.Object, err error) {
	if ast.Kind != Root {
		return nil, errors.New("available for root node only")
	}

	for _, st := range ast.Subtrees {
		res, err = st.eval()
		if err != nil {
			return nil, err
		}
	}

	return
}

func (ast *AST) eval() (types.Object, error) {
	if ast.Kind == CallExpr {
		return ast.call()
	}

	if ast.Kind == VariableRef {
		return ast.getVar()
	}

	if ast.Kind == Literal {
		switch ast.Token.Type() {
		case String:
			return types.String(ast.Token.Value()), nil
		case Int:
			num, err := strconv.ParseInt(ast.Token.Value(), 10, 64)
			if err != nil {
				return nil, err
			}

			return types.NumberFrom(num), nil
		case Float:
			num, err := strconv.ParseFloat(ast.Token.Value(), 64)
			if err != nil {
				return nil, err
			}

			return types.NumberFrom(num), nil
		default:
		}
	}

	return nil, fmt.Errorf("unimplemented")
}

func (ast *AST) getVar() (types.Object, error) {
	def, ok := ast.Ctx.FindDef(ast.Identifier())
	if !ok {
		return nil, fmt.Errorf(`"%v" is not defined`, ast.Identifier())
	}

	return def, nil
}

func (ast *AST) call() (types.Object, error) {
	def, ok := ast.Ctx.FindDef(ast.Identifier())
	if !ok {
		return nil, fmt.Errorf(`"%v" is not defined`, ast.Identifier())
	}

	fun, ok := def.(types.Callable)
	if !ok {
		return nil, fmt.Errorf(`"%v" is not a function`, ast.Identifier())
	}

	var args []types.Object
	for _, st := range ast.Subtrees {
		arg, err := st.eval()
		if err != nil {
			return nil, err
		}

		args = append(args, arg)
	}

	return fun.Call(args...)
}
