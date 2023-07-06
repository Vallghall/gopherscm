package types

import (
	"fmt"
	"github.com/Vallghall/gopherscm/internal/core/operator"
	"github.com/Vallghall/gopherscm/internal/errscm"
)

// number - supported number types
type number uint

// number enum
const (
	Int number = iota
	Float
)

// Number - wrapper around operations defined on numbers
type Number struct {
	t     number
	value any
}

func (n *Number) String() string {
	return fmt.Sprintf("%v", n.value)
}

// Value - return the value of number.
// Implements data.Object interface
func (n *Number) Value() any {
	return n.value
}

// NewNumber - Number type constructor
func NewNumber(d int64) *Number {
	return &Number{
		t:     Int,
		value: d,
	}
}

// NumberFrom - Number type constructor from a number
func NumberFrom(obj any) *Number {
	var t number
	switch obj.(type) {
	case int64:
		t = Int
	case float64:
		t = Float
	}

	return &Number{
		t:     t,
		value: obj,
	}
}

func (n *Number) Int() int64 {
	return n.value.(int64)
}

func (n *Number) Float() float64 {
	return n.value.(float64)
}

// ApplyOperation - applies operation considering underlying types
func (n *Number) ApplyOperation(op operator.Operator, o Object) (obj *Number, err error) {
	num, ok := o.(*Number)
	if !ok {
		return nil, errscm.ErrNaN
	}

	switch num.value.(type) {
	case int64:
		obj = n.ApplyInt(op, num)
	case float64:
		obj = n.ApplyFloat(op, num)
	default:
		return nil, errscm.ErrUnsupported
	}

	return
}

// ApplyUnary - applies unary operator considering underlying types
// Right now it is implemented for negation only
func (n *Number) ApplyUnary() (obj Object, err error) {
	switch n.t {
	case Int:
		obj = NumberFrom(operator.Neg(n.Int()))
	case Float:
		obj = NumberFrom(operator.Neg(n.Float()))
	default:
		return nil, errscm.ErrUnsupported
	}

	return
}

// ApplyInt - handles int64 application to Number
func (n *Number) ApplyInt(op operator.Operator, num *Number) (obj *Number) {
	if n.t == Int {
		obj = NumberFrom(operator.Run(op, n.Int(), num.Int()))
	} else {
		obj = NumberFrom(operator.Run(op, n.Float(), float64(num.Int())))
	}

	return
}

// ApplyFloat - handles float64 application to Number
func (n *Number) ApplyFloat(op operator.Operator, num *Number) (obj *Number) {
	if n.t == Float {
		obj = NumberFrom(operator.Run(op, n.Float(), num.Float()))
	} else {
		obj = NumberFrom(operator.Run(op, float64(n.Int()), num.Float()))
	}

	return
}
