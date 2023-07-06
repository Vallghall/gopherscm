package parser

import (
	"github.com/Vallghall/gopherscm/internal/data"
)

const (
	lParen = "("
	rParen = ")"
)

// Parse - parsing token stream into AST
func Parse(ts data.TokenStream) *data.AST {
	ast := data.ASTRoot()
	idx := 0
	_ = parse(ast, ts, idx)
	return ast
}

// parse - recursive helper called from Parse
func parse(ast *data.AST, ts data.TokenStream, idx int) int {
	for idx < len(ts) {
		token := ts[idx]
		if token.Type() == data.Syntax {
			if token.Value() == lParen {
				idx++
				subtree := ast.Nest(ts[idx])
				idx = parse(subtree, ts, idx+1)
				continue
			}

			if token.Value() == rParen {
				return idx + 1
			}

		}

		idx++
		ast.Add(token)
	}

	return idx
}
