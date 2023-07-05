package data

// Context - context of the scope
type Context struct {
	// outerCtx - all the scope's outer scopes til the global context
	outerCtx *Context
	// probably replace with sync.Map
	symbolTable map[string]any
}

// NewContext - Context constructor for the global
// context that has no outer context
func NewContext(t map[string]any) *Context {
	return &Context{
		symbolTable: t,
	}
}

// Spawn - creates child context, that points to
// the parent context, with empty symbol table
func (c *Context) Spawn() *Context {
	return &Context{
		outerCtx:    c,
		symbolTable: make(map[string]any),
	}
}

// FindDef - find definition for the given identifier
func (c *Context) FindDef(s string) (any, bool) {
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
