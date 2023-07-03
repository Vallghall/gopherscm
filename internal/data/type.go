package data

// Type - supported kinds of token
// TODO: support all the other data
type Type uint

// Type enum
const (
	Syntax Type = iota
	Id
	Int
	String
)
