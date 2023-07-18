package interp

import (
	"fmt"

	"github.com/Vallghall/gopherscm/internal/core/types"
	"github.com/Vallghall/gopherscm/internal/data"
)

// Func - ast, evaluated as function body
type Func struct {
	*data.AST
	Params []string
}

// Func - create function from AST subtrees
func NewFunc(ast *data.AST, body []*data.AST) *Func {
	return &Func{
		AST: &data.AST{
			Ctx:      ast.Ctx.Spawn(),
			Kind:     data.Function,
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
		result, err = Eval(expr)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}
