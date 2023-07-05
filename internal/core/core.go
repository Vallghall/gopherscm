package core

import (
	"github.com/Vallghall/gopherscm/internal/core/arithmetics"
	"github.com/Vallghall/gopherscm/internal/core/stdio"
	"github.com/Vallghall/gopherscm/internal/core/types"
)

// DefaultDefinitions - returns a symbol table with builtin definitions
func DefaultDefinitions() map[string]types.Object {
	return map[string]types.Object{
		// Arithmetics
		"+": arithmetics.Primitive(arithmetics.Plus),
		"-": arithmetics.Primitive(arithmetics.Minus),
		"*": arithmetics.Primitive(arithmetics.Multiply),
		"/": arithmetics.Primitive(arithmetics.Divide),

		// Standart output
		"display":   stdio.IOHandler(stdio.Display),
		"newline":   stdio.IOHandler(stdio.NewLine),
		"displayln": stdio.IOHandler(stdio.Displayln),
	}
}
