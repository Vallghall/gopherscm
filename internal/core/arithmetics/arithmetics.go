package arithmetics

import (
	"fmt"
	"github.com/Vallghall/gopherscm/internal/core/operator"
	"github.com/Vallghall/gopherscm/internal/core/types"
	"github.com/Vallghall/gopherscm/internal/errscm"
)

// Primitive - wrapper for primitive arithmetics
type Primitive func(args ...types.Object) (types.Object, error)

func (p Primitive) Call(args ...types.Object) (types.Object, error) {
	return p(args...)
}

func (p Primitive) Value() any {
	return "PrimitiveOperation"
}

// Plus – `+` primitive
func Plus(args ...types.Object) (types.Object, error) {
	sum := types.NewNumber(0)
	var err error
	for _, arg := range args {
		sum, err = sum.ApplyOperation(operator.Addition, arg)
		if err != nil {
			return nil, err
		}
	}

	return sum, nil
}

// Minus – `-` primitive
func Minus(args ...types.Object) (types.Object, error) {
	if len(args) < 1 {
		return nil, errscm.ErrTooLittleArguments
	}

	var err error
	sub, ok := args[0].(*types.Number)
	if !ok {
		return nil, errscm.ErrNaN
	}

	if len(args) == 1 {
		return sub.ApplyUnary()
	}

	for _, arg := range args[1:] {
		sub, err = sub.ApplyOperation(operator.Subtraction, arg)
		if err != nil {
			return nil, err
		}
	}

	return sub, nil
}

// Multiply – `*` primitive
func Multiply(args ...types.Object) (types.Object, error) {
	prod := types.NewNumber(1)
	var err error
	for _, arg := range args {
		prod, err = prod.ApplyOperation(operator.Multiplication, arg)
		if err != nil {
			return nil, err
		}
	}

	return prod, nil
}

// Divide – `/` primitive
func Divide(args ...types.Object) (types.Object, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("%w: expected 2 arguments, got: %d", errscm.ErrUnexpectedNumberOfArguments, len(args))
	}

	a, ok := args[0].(*types.Number)
	if !ok {
		return nil, errscm.ErrNaN
	}
	b := args[1]
	result, err := a.ApplyOperation(operator.Division, b)
	if err != nil {
		return nil, err
	}

	return result, nil
}
