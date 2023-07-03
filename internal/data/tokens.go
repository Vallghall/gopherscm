package data

// TokenStream - alias
type TokenStream []*Token

// Token representing a Scheme lexical unit
type Token struct {
	value string
	t     Type
	meta  *Meta
}

// Value - Token value getter
func (tp *Token) Value() string {
	return tp.value
}

// Type - token t (type) getter
func (tp *Token) Type() Type {
	return tp.t
}

// Meta - token meta getter
func (tp *Token) Meta() *Meta {
	return tp.meta
}

// NewToken - token constructor from value and token type
func NewToken(value string, t Type) *Token {
	return &Token{
		value: value,
		t:     t,
	}
}

// TokenFromMeta - token constructor that only sets meta
func TokenFromMeta(m *Meta) *Token {
	return &Token{
		meta: m.Current(),
	}
}

// Set - sets the type and value of a token
func (tp *Token) Set(t Type, syms ...rune) *Token {
	tp.t, tp.value = t, string(syms)
	return tp
}
