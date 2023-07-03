package data

// AST - represents Scheme program structure
type AST struct {
	OuterCtx *Context
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

// Add - AST node constructor
func (ast *AST) Add(t *Token) *AST {
	node := &AST{
		OuterCtx: ast.Ctx,
		Ctx:      ast.Ctx,
		Token:    t,
		//kind:  TODO: add kind parsing,
		Subtrees: make([]*AST, 0),
	}

	ast.Subtrees = append(ast.Subtrees, node)

	return node
}

// Expr - expressions supported by the program
type Expr uint

// Expr enum
const (
	// Literal - expression which is evaluated into itself
	Literal Expr = iota
	// CallExpr - function with a list of arguments
	CallExpr
	// DefineExpr - expression that creates a new lexical scope
	// and binds variable with an expression
	DefineExpr

	// Root - AST root unique expressions kind
	Root = 9999
)
