package lexer

import (
	"errors"
	"unicode"
)

// Lex transforms input slice of symbols into slice of Scheme tokens
func Lex(src []rune) (TokenStream, error) {
	ts := make(TokenStream, 0)
	m := newMeta()
	inputLength := len(src)
	cursor := 0

	var token *Token
	var err error
	parenCount := 0

	// cursor is updated in Tokenize and skip* functions only
	for cursor < inputLength {
		old := cursor
		cursor = skipSingleLineComment(cursor, src, m)
		cursor = skipWhiteSpaces(cursor, src, m)
		if old != cursor {
			continue // recheck whitespaces and comments after theit deletion
		}

		cursor, token, err = Tokenize(cursor, src, m)
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
	}

	if parenCount > 0 {
		return nil, ErrMissingClosingParenthesis
	}

	return ts, nil
}

// Tokenize - Extracts token from rune sequence
func Tokenize(cursor int, src []rune, m *meta) (int, *Token, error) {
	if cursor >= len(src) {
		return cursor, nil, ErrEndOfInput
	}

	sym := src[cursor]

	// '(' and ')' are the only syntax tokens
	if sym == '(' || sym == ')' {
		m.inc()
		return cursor + 1, NewToken(Syntax, m, sym), nil
	}

	// parsing string literal like "foo"
	if sym == '"' {
		return extractString(cursor, src, m)
	}

	// parsing integer literal
	// TODO: add floating point numbers lexing
	if unicode.IsDigit(sym) {
		return extractNumber(cursor, src, m)
	}

	// Check for an identifier
	if isValidChar(sym) {
		return extractIdentifier(cursor, src, m)
	}

	return cursor, nil, ErrInvalidSymbol
}

// skipSingleLineComment
func skipSingleLineComment(cursor int, src []rune, m *meta) int {
	inputLength := len(src)
	if src[cursor] == ';' {
		for src[cursor] != '\n' {
			m.inc()
			cursor++
			if cursor >= inputLength {
				return cursor
			}
		}
		m.newLine()
		cursor++
	}

	return cursor
}

// skipWhitespaces - helper func for ommiting whitespaces
func skipWhiteSpaces(cursor int, src []rune, m *meta) int {
	inputLength := len(src)
	for unicode.IsSpace(src[cursor]) {
		cursor++
		if cursor >= inputLength {
			return cursor
		}

		m.incNL(src[cursor])
	}

	return cursor
}

// extractString - helper func for extracting String token
// FIXME: fix so that only single-line strings are allowed
// TODO: add support for escape sequences
func extractString(cursor int, src []rune, m *meta) (int, *Token, error) {
	cursor++ // move forward from quote
	if cursor >= len(src) {
		return cursor, nil, ErrEndOfInput
	}
	m.inc()

	str := make([]rune, 0)
	for sym := src[cursor]; sym != '"'; sym = src[cursor] {
		str = append(str, sym)
		cursor++
		if cursor >= len(src) {
			return cursor, nil, ErrMissingMatchingDoubleQuotes
		}

		m.inc()
	}

	m.inc()
	cursor++
	return cursor, NewToken(String, m, str...), nil
}

// extractNumber - helper func for extracting an Int token
// TODO: add floating point token support
func extractNumber(cursor int, src []rune, m *meta) (int, *Token, error) {
	number := []rune{src[cursor]}
	cursor++
	if cursor >= len(src) {
		return cursor, nil, ErrEndOfInput
	}
	m.inc()

	for unicode.IsDigit(src[cursor]) {
		number = append(number, src[cursor])
		cursor++
		if cursor >= len(src) {
			return cursor, nil, ErrEndOfInput
		}

		m.inc()
	}

	cursor = skipSingleLineComment(cursor, src, m)
	if cursor >= len(src) {
		return cursor, nil, ErrEndOfInput
	}

	if sym := src[cursor]; !(unicode.IsSpace(sym) || sym == ')') {
		return cursor, nil, ErrInvalidIntegerLiteral
	}

	return cursor, NewToken(Int, m, number...), nil
}

// extractIdentifier - helper func for lexing identifiers
func extractIdentifier(cursor int, src []rune, m *meta) (int, *Token, error) {
	id := []rune{src[cursor]}
	cursor++
	if cursor >= len(src) {
		return cursor, nil, ErrEndOfInput
	}
	m.incNL(src[cursor])

	for isValidChar(src[cursor]) || unicode.IsDigit(src[cursor]) {
		id = append(id, src[cursor])

		cursor++
		if cursor >= len(src) {
			return cursor, nil, ErrEndOfInput
		}
		m.incNL(src[cursor])
	}

	return cursor, NewToken(Id, m, id...), nil
}

// isValidChar - predicate for checking a valid identifier's symbol
func isValidChar(sym rune) bool {
	return unicode.IsLetter(sym) || sym == '?' || sym == '!' || sym == '-' || sym == '_' || sym == '+'
}
