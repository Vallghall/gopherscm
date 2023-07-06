package data

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Vallghall/gopherscm/internal/errscm"

	"github.com/Vallghall/gopherscm/internal/core"
	"github.com/Vallghall/gopherscm/internal/core/types"
)

// AST - represents Scheme program structure
type AST struct {
	Ctx      *Context
	Token    *Token
	Kind     Expr
	Subtrees []*AST
}

// Func - ast, evaluated as function body
type Func struct {
	*AST
	Params []string
}

// Func - create function from AST subtrees
func (ast *AST) Func(body []*AST) *Func {
	return &Func{
		AST: &AST{
			Ctx:      ast.Ctx.Spawn(),
			Kind:     Function,
			Subtrees: body,
		},
	}
}

// Value - types.Object interface implementation
func (f *Func) Value() any {
	return "lambda" // TODO: improve it
}

// Call - types.Callable interface implementation
// Binds given arguments to parameter list and evaluates the Func
func (f *Func) Call(args ...types.Object) (result types.Object, err error) {
	if len(args) != len(f.Params) {
		return nil, fmt.Errorf(
			"not enough arguments:\nexpected %d\ngot: %d",
			len(args),
			len(f.Params),
		)
	}

	// define function params within the current context
	for i, key := range f.Params {
		f.Ctx.Set(key, args[i])
	}

	for _, expr := range f.Subtrees {
		expr.Ctx = f.Ctx // enforce function ctx onto its body expressions
		result, err = expr.eval()
		if err != nil {
			return nil, err
		}
	}

	return result, nil
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

// eval - evaluates expression based on its kind
func (ast *AST) eval() (types.Object, error) {
	switch ast.Kind {
	case CallExpr:
		return ast.call()
	case VariableRef:
		return ast.getVar()
	case DefineExpr:
		return ast.define()
	case Literal:
		return ast.evalLiteral()
	}

	return nil, fmt.Errorf("unimplemented")
}

// evalLiteral - wrapping literal's token value into
// a type that implements types.Object
func (ast *AST) evalLiteral() (types.Object, error) {
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

	return nil, errors.New("unimplemented")
}

// getVar - variable lookup
func (ast *AST) getVar() (types.Object, error) {
	def, ok := ast.Ctx.FindDef(ast.Identifier())
	if !ok {
		return nil, fmt.Errorf(`"%v" is not defined`, ast.Identifier())
	}

	return def, nil
}

// call - asserts that the object called is types.Callable,
// evaluates its list of arguments and calls the function with
// the evaluated arguments
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
		st.Ctx = ast.Ctx // function context enforcement
		arg, err := st.eval()
		if err != nil {
			return nil, err
		}

		args = append(args, arg)
	}

	return fun.Call(args...)
}

// define - handles variable and function definitions
func (ast *AST) define() (types.Object, error) {
	id := ast.Subtrees[0]

	if id.Kind == VariableRef {
		if len(ast.Subtrees) != 2 {
			return nil, fmt.Errorf("%w: expected 2 args, got: %d", errscm.ErrUnexpectedNumberOfArguments, len(ast.Subtrees))
		}

		def := ast.Subtrees[1]
		value, err := def.eval()
		if err != nil {
			return nil, err
		}

		ast.Ctx.Set(id.Identifier(), value)
		return nil, nil
	}

	if id.Kind == CallExpr {
		if len(ast.Subtrees) < 2 {
			return nil, fmt.Errorf("%w: missing function body", errscm.ErrTooLittleArguments)
		}

		fn := id.Func(ast.Subtrees[1:])
		for _, param := range id.Subtrees {
			if param.Kind == VariableRef {
				fn.Params = append(fn.Params, param.Identifier())
				continue
			}

			return nil, fmt.Errorf("%s is not a valid identifier", param.Identifier())
		}

		ast.Ctx.Set(id.Identifier(), fn)
	}

	return nil, nil
}
