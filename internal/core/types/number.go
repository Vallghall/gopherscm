package types

import (
	"errors"
	"fmt"
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

// Value - return the value of number.
// Implements data.Object interface
func (n *Number) Value() any {
	return n.value
}

// NewNumber - Number type constructor
func NewNumber() *Number {
	return &Number{
		t:     Int,
		value: int64(0),
	}
}

// NumberFrom - Number type constructor from a number
func NumberFrom(obj any, kind number) *Number {
	return &Number{
		t:     kind,
		value: obj,
	}
}

func (n *Number) Int() int64 {
	return n.value.(int64)
}

func (n *Number) Float() float64 {
	return n.value.(float64)
}

// Add - addition that considers the underlying type of Number
func (n *Number) Add(obj Object) error {
	num, ok := obj.(*Number)
	if !ok {
		return errors.New("Not a Number")
	}

	switch num.value.(type) {
	case int64:
		n.AddInt(num)
	case float64:
		n.AddFloat(num)
	default:
		return fmt.Errorf("unsupported")
	}

	return nil
}

// AddInt - handles int64 addition to Number
func (n *Number) AddInt(num *Number) {
	if n.t == Int {
		n.value = n.Int() + num.Int()
	} else {
		n.value = n.Float() + float64(num.Int())
	}
}

// AddFloat - handles float64 addition to Number
func (n *Number) AddFloat(num *Number) {
	if n.t == Float {
		n.value = n.Float() + num.Float()
	} else {
		// if t = Int, upgrade it to Float
		n.t = Float
		n.value = float64(n.Int()) + num.Float()
	}
}

func (n *Number) String() string {
	return fmt.Sprintf("%v", n.value)
}
