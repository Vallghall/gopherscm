package types

type String string

func (s String) Value() any {
	return s
}
