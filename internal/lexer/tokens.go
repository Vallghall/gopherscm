package lexer

// TokenStream - alias
type TokenStream []*Token

// Token representing a Scheme lexical unit
type Token struct {
	value string
	t     Type
	// token location
	line int
	pos  int
}

// NewToken - token constructor
func NewToken(t Type, m *meta, syms ...rune) *Token {
	return &Token{
		value: string(syms),
		t:     t,
		line:  m.line,
		pos:   m.pos,
	}
}

// Type - supported kinds of token
// TODO: support all the other types
type Type uint

// Type enum
const (
	Syntax Type = iota
	Id
	Int
	String
)
