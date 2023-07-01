package lexer

import "errors"

var (
	ErrCursorOutOfRange            = errors.New("cursor is out of range")
	ErrMissingMatchingDoubleQuotes = errors.New("missing matching double quotes")
	ErrInvalidIntegerLiteral       = errors.New("invalid integer literal")
	ErrInvalidSymbol               = errors.New("invalid symbol")
	ErrFreeClosingParantesis       = errors.New("free closing parenthesis")
	ErrMissingClosingParenthesis   = errors.New("missing matching closing parenthesis")
)
