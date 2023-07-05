package data

// AST - represents Scheme program structure
type AST struct {
	Ctx      *Context
	Token    *Token
	Kind     Expr
	Subtrees []*AST
}

// ASTRoot - constructor for AST
func ASTRoot() *AST {
	return &AST{
		Kind:     Root,
		Subtrees: make([]*AST, 0),
	}
}

// Nest - AST node constructor
func (ast *AST) Nest(t *Token) *AST {
	node := &AST{
		Ctx:      ast.Ctx.Spawn(),
		Token:    t,
		Kind:     CallExpr, // all nested forms have functions as the first elem
		Subtrees: make([]*AST, 0),
	}

	ast.Subtrees = append(ast.Subtrees, node)

	return node
}

// Add - AST node constructor
func (ast *AST) Add(t *Token) *AST {
	var e Expr
	switch t.Type() {
	case Int, Float, String:
		e = Literal
	case Id:
		e = VariableRef
	default: // fill in later
	}

	node := &AST{
		Ctx:      ast.Ctx,
		Token:    t,
		Kind:     e,
		Subtrees: make([]*AST, 0),
	}

	ast.Subtrees = append(ast.Subtrees, node)

	return node
}

// Eval - evaluate AST to a sigle value
func (ast *AST) Eval() any {
	for _, _ = range ast.Subtrees {
		//TODO
	}
	return nil
}

func (ast *AST) eval() any {
	return nil
}
