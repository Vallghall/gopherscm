package types

// Callable - function interface
type Callable interface {
	Call(args ...Object) (Object, error)
}

// Object - any Scheme value
type Object interface {
	Value() any
}
