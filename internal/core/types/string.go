package types

// String - wrapper for strings that implements Object
type String string

// Value - Object implementation
func (s String) Value() any {
	return s
}
