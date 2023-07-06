package stdio

import (
	"fmt"
	"github.com/Vallghall/gopherscm/internal/errscm"

	"github.com/Vallghall/gopherscm/internal/core/types"
)

// IOHandler - wrapper around buitin funcs related to io operations
// defined in the standart
type IOHandler func(args ...types.Object) (types.Object, error)

func (io IOHandler) Call(args ...types.Object) (types.Object, error) {
	return io(args...)
}

func (p IOHandler) Value() any {
	return "PrimitiveOperation"
}

// Display - prints given args to stdout
func Display(args ...types.Object) (types.Object, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("%w: expected 1 arg, got %d", errscm.ErrUnexpectedNumberOfArguments, len(args))
	}

	fmt.Print(args[0])
	return nil, nil
}

// NewLine - prints new line character to stdout
// args are to match the IOHandler alias, they're not used
func NewLine(args ...types.Object) (types.Object, error) {
	fmt.Println()
	return nil, nil
}

// Displayln - prints given args to stdout and adds
// a new line character at the end
func Displayln(args ...types.Object) (types.Object, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("%w: expected 1 arg, got %d", errscm.ErrUnexpectedNumberOfArguments, len(args))
	}

	fmt.Println(args[0])
	return nil, nil
}
