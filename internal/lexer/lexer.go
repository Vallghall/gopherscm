package lexer

import (
	"errors"
	"unicode"
)

// Lex transforms input slice of symbols into slice of Scheme tokens
func Lex(src []rune) ([]Token, error) {
	ts := make([]Token, 0)

	inputLength := len(src)
	cursor := 0

	skipWhiteSpaces := func() {
		if cursor < inputLength {
			for unicode.IsSpace(src[cursor]) {
				cursor++
			}
		}
	}

	var token Token
	var err error
	parenCount := 0
	for {
		skipWhiteSpaces()

		cursor, token, err = Tokenize(src, cursor)
		if err != nil {
			if errors.Is(err, ErrEndOfInput) {
				break
			}
			return nil, err
		}

		ts = append(ts, token)
		if token.t == Syntax {
			if token.value == "(" {
				parenCount++
			} else {
				parenCount--
			}
		}

		if parenCount < 0 {
			return nil, ErrFreeClosingParantesis
		}

		if cursor >= inputLength {
			break
		}
	}

	if parenCount > 0 {
		return nil, ErrMissingClosingParenthesis
	}

	return ts, nil
}

// Tokenize - Extracts token from rune sequence
func Tokenize(src []rune, cursor int) (int, Token, error) {
	if cursor >= len(src) {
		return cursor, Token{}, ErrEndOfInput
	}

	sym := src[cursor]

	// '(' and ')' are the only syntax tokens
	if sym == '(' || sym == ')' {
		return cursor + 1, Token{value: string(sym), t: Syntax}, nil
	}

	// parsing string literal like "foo"
	if sym == '"' {
		return extractString(cursor, src)
	}

	// parsing integer literal
	// TODO: add floating point numbers lexing
	if unicode.IsDigit(sym) {
		return extractNumber(cursor, src)
	}

	// Check for an identifier
	if isValidChar(sym) {
		return extractIdentifier(cursor, src)
	}

	return cursor, Token{}, ErrInvalidSymbol
}

// extractString - helper func for extracting String token
func extractString(cursor int, src []rune) (int, Token, error) {
	cursor++ // move forward from quote
	if cursor >= len(src) {
		return cursor, Token{}, ErrEndOfInput
	}

	str := make([]rune, 0)
	for sym := src[cursor]; sym != '"'; sym = src[cursor] {
		str = append(str, sym)
		cursor++
		if cursor >= len(src) {
			return cursor, Token{}, ErrMissingMatchingDoubleQuotes
		}
	}

	cursor++
	return cursor, Token{value: string(str), t: String}, nil
}

// extractNumber - helper func for extracting an Int token
// TODO: add floating point token support
func extractNumber(cursor int, src []rune) (int, Token, error) {
	number := []rune{src[cursor]}
	cursor++
	if cursor >= len(src) {
		return cursor, Token{}, ErrEndOfInput
	}

	for unicode.IsDigit(src[cursor]) {
		number = append(number, src[cursor])
		cursor++

		if cursor >= len(src) {
			return cursor, Token{}, ErrEndOfInput
		}
	}

	if sym := src[cursor]; !(unicode.IsSpace(sym) || sym == ')') {
		return cursor, Token{}, ErrInvalidIntegerLiteral
	}

	return cursor, Token{value: string(number), t: Int}, nil
}

// extractIdentifier - helper func for lexing identifiers
func extractIdentifier(cursor int, src []rune) (int, Token, error) {
	id := []rune{src[cursor]}
	cursor++
	if cursor >= len(src) {
		return cursor, Token{}, ErrEndOfInput
	}

	for isValidChar(src[cursor]) || unicode.IsDigit(src[cursor]) {
		id = append(id, src[cursor])

		cursor++
		if cursor >= len(src) {
			return cursor, Token{}, ErrEndOfInput
		}
	}

	return cursor, Token{value: string(id), t: Id}, nil
}

// isValidChar - predicate for checking a valid identifier's symbol
func isValidChar(sym rune) bool {
	return unicode.IsLetter(sym) || sym == '?' || sym == '!' || sym == '-' || sym == '_' || sym == '+'
}
