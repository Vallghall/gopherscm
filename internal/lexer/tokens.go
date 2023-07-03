package lexer

import "github.com/Vallghall/gopherscm/internal/types"

// TokenStream - alias
type TokenStream []*Token

// Token representing a Scheme lexical unit
type Token struct {
	value string
	t     Type
	meta  *types.Meta
}

func (tp *Token) Value() string {
	return tp.value
}

func (tp *Token) Type() Type {
	return tp.t
}

func (tp *Token) Meta() *types.Meta {
	return tp.meta
}

/*
// newToken - token constructor
func newToken(t Type, m *types.Meta, syms ...rune) *Token {
	return &Token{
		value: string(syms),
		t:     t,
		meta:  m,
	}
}
*/

// tokenFromMeta - token constructor that only sets meta
func tokenFromMeta(m *types.Meta) *Token {
	return &Token{
		meta: m.Current(),
	}
}

// Set - sets the type and value of a token
func (tp *Token) Set(t Type, syms ...rune) *Token {
	tp.t, tp.value = t, string(syms)
	return tp
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
