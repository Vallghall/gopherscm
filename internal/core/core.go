package core

import (
	"github.com/Vallghall/gopherscm/internal/core/arithmetics"
	"github.com/Vallghall/gopherscm/internal/core/types"
)

func DefaultDefinitions() map[string]types.Object {
	return map[string]types.Object{
		"+": arithmetics.Primitive(arithmetics.Plus),
	}
}
