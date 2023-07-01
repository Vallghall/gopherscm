package lexer

type Token struct {
	value string
	t     Type
}

type Type uint

const (
	Syntax Type = iota
	Id
	Int
	String
)
