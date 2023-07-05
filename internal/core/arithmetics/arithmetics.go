package arithmetics

import (
	"github.com/Vallghall/gopherscm/internal/core/types"
)

// Primitive - wrapper for primitive arithmetics
type Primitive func(args ...types.Object) (types.Object, error)

func (p Primitive) Call(args ...types.Object) (types.Object, error) {
	return p(args...)
}

func (p Primitive) Value() any {
	return "PrimitiveOperation"
}

// Plus – + primitive
func Plus(args ...types.Object) (types.Object, error) {
	sum := types.NewNumber()
	for _, arg := range args {
		err := sum.Add(arg)
		if err != nil {
			return nil, err
		}
	}

	return sum, nil
}
