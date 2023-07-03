package parser

import (
	"github.com/Vallghall/gopherscm/internal/data"
)

// TODO: make parser; make ast interface in
// a new package for decoupling from parser package

// ast - abstract syntax tree.
// Represents the program structure
type ast struct {
	car *ast   // tree root node
	cdr []*ast // node leaves
}

// Parse - parsing token stream into AST
func Parse(ts []data.Token) *ast {
	return nil // TODO: implement!
}
