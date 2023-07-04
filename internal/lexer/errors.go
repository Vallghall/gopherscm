package lexer

import (
	"errors"
	"fmt"
)

// SyntaxError - error type that includes
// useful information about error positioning
type SyntaxError struct {
	err      error
	line     int
	position int
}

// ReportError - constructor for SyntaxError
func ReportError(line, pos int, err error) error {
	return &SyntaxError{
		err:      err,
		line:     line,
		position: pos,
	}
}

// Error - error interface implementation
func (se *SyntaxError) Error() string {
	return fmt.Sprintf("ERROR: Syntax error at line %d, position %d: %s", se.line, se.position, se.err.Error())
}

// Unwrap - unwrap interface implementation
func (se *SyntaxError) Unwrap() error {
	return se.err
}

var (
	// ErrEndOfInput - signals of unexpected end of input
	ErrEndOfInput                  = errors.New("cursor is out of range")
	ErrMissingMatchingDoubleQuotes = errors.New("missing matching double quotes")
	ErrInvalidIntegerLiteral       = errors.New("invalid integer literal")
	ErrInvalidSymbol               = errors.New("invalid symbol")
	ErrFreeClosingParantesis       = errors.New("free closing parenthesis")
	ErrMissingClosingParenthesis   = errors.New("missing matching closing parenthesis")
	ErrNaN                         = errors.New("NaN")
)
