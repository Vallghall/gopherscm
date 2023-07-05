package data

import "github.com/Vallghall/gopherscm/internal/core/types"

// Context - context of the scope
type Context struct {
	// outerCtx - all the scope's outer scopes til the global context
	outerCtx *Context
	// probably replace with sync.Map
	symbolTable map[string]types.Object
}

// NewContext - Context constructor for the global
// context that has no outer context
func NewContext(t map[string]types.Object) *Context {
	return &Context{
		symbolTable: t,
	}
}

// Spawn - creates child context, that points to
// the parent context, with empty symbol table
func (c *Context) Spawn() *Context {
	return &Context{
		outerCtx:    c,
		symbolTable: make(map[string]types.Object),
	}
}

// FindDef - find definition for the given identifier
func (c *Context) FindDef(s string) (types.Object, bool) {
	def, ok := c.symbolTable[s]
	if !ok {
		// only global context has nil value
		if c.outerCtx != nil {
			return c.outerCtx.FindDef(s)
		}

		return nil, false
	}

	return def, true
}

// Set - binds provided object to a given key within the current context
func (c *Context) Set(key string, obj types.Object) {
	c.symbolTable[key] = obj
}
