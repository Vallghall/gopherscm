package interp

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Vallghall/gopherscm/internal/core/types"
	"github.com/Vallghall/gopherscm/internal/data"
	"github.com/Vallghall/gopherscm/internal/errscm"
)

// Walk - evaluate whole AST to a single value
func Walk(ast *data.AST) (res types.Object, err error) {
	if ast.Kind != data.Root {
		return nil, errors.New("available for root node only")
	}

	for _, st := range ast.Subtrees {
		res, err = Eval(st)
		if err != nil {
			return nil, err
		}
	}

	return
}

// Eval - evaluates expression subtree based on its kind
func Eval(ast *data.AST) (types.Object, error) {
	switch ast.Kind {
	case data.CallExpr:
		return call(ast)
	case data.VariableRef:
		return getVar(ast)
	case data.DefineExpr:
		return define(ast)
	case data.Literal:
		return evalLiteral(ast)
	}

	return nil, errscm.ErrUnsupported
}

// evalLiteral - wrapping literal's token value into
// a type that implements types.Object
func evalLiteral(ast *data.AST) (types.Object, error) {
	switch ast.Token.Type() {
	case data.String:
		return types.String(ast.Token.Value()), nil
	case data.Int:
		num, err := strconv.ParseInt(ast.Token.Value(), 10, 64)
		if err != nil {
			return nil, err
		}

		return types.NumberFrom(num), nil
	case data.Float:
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
func getVar(ast *data.AST) (types.Object, error) {
	def, ok := ast.Ctx.FindDef(ast.Identifier())
	if !ok {
		return nil, fmt.Errorf(`"%v" is not defined`, ast.Identifier())
	}

	return def, nil
}

// call - asserts that the object called is types.Callable,
// evaluates its list of arguments and calls the function with
// the evaluated arguments
func call(ast *data.AST) (types.Object, error) {
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
		arg, err := Eval(st)
		if err != nil {
			return nil, err
		}

		args = append(args, arg)
	}

	return fun.Call(args...)
}

// define - handles variable and function definitions
func define(ast *data.AST) (types.Object, error) {
	id := ast.Subtrees[0]

	if id.Kind == data.VariableRef {
		if len(ast.Subtrees) != 2 {
			return nil, fmt.Errorf("%w: expected 2 args, got: %d", errscm.ErrUnexpectedNumberOfArguments, len(ast.Subtrees))
		}

		def := ast.Subtrees[1]
		value, err := Eval(def)
		if err != nil {
			return nil, err
		}

		ast.Ctx.Set(id.Identifier(), value)
		return nil, nil
	}

	if id.Kind == data.CallExpr {
		if len(ast.Subtrees) < 2 {
			return nil, fmt.Errorf("%w: missing function body", errscm.ErrTooLittleArguments)
		}

		fn := NewFunc(id, ast.Subtrees[1:])
		for _, param := range id.Subtrees {
			if param.Kind == data.VariableRef {
				fn.Params = append(fn.Params, param.Identifier())
				continue
			}

			return nil, fmt.Errorf("%s is not a valid identifier", param.Identifier())
		}

		ast.Ctx.Set(id.Identifier(), fn)
	}

	return nil, nil
}
