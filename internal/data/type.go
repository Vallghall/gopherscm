package data

import (
	"encoding/json"
	"errors"
)

var ErrUnsupportedTokenType = errors.New("unsupported token type")

// Type - supported kinds of token
// TODO: support all the other data
type Type uint

// Type enum
const (
	Syntax Type = iota
	Id
	Int
	Float
	String
)

func (t Type) MarshalJSON() ([]byte, error) {
	switch t {
	case Syntax:
		return json.Marshal("Syntax")
	case Id:
		return json.Marshal("Identifier")
	case Int:
		return json.Marshal("Integer")
	case String:
		return json.Marshal("String")
	default:
		return nil, ErrUnsupportedTokenType
	}
}
