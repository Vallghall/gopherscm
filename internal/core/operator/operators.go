package operator

// Operator - list of defined operators
type Operator uint

// Operator enum
const (
	Addition Operator = iota
	Subtraction
	Multiplication
	Division
)

// Run - executes given operation for a generic number type
func Run[T int64 | float64 | complex128](op Operator, a T, b T) (result T) {
	switch op {
	case Addition:
		return a + b
	case Subtraction:
		return a - b
	case Multiplication:
		return a * b
	case Division:
		return a / b
	}

	return
}

// Neg - negates given generic number
func Neg[T int64 | float64 | complex128](a T) T {
	return -a
}
