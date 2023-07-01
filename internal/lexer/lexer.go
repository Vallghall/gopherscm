package lexer

import (
	"unicode"
)

// Tokenize - Extracts token from rune sequence
func Tokenize(src []rune, cursor int) (int, Token, error) {
	if cursor >= len(src) {
		return cursor, Token{}, ErrCursorOutOfRange
	}

	sym := src[cursor]
	// '(' and ')' are the only syntax tokens
	if sym == '(' || sym == ')' {
		return cursor + 1, Token{value: string(sym), t: Syntax}, nil
	}

	// parsing string literal like "foo"
	if sym == '"' {
		cursor++
		str := make([]rune, 0)
		for sym = src[cursor]; sym != '"'; sym = src[cursor] {
			str = append(str, sym)
			cursor++
			if cursor >= len(src) {
				return cursor, Token{}, ErrMissingMatchingDoubleQuotes
			}
		}

		cursor++
		return cursor, Token{value: string(str), t: String}, nil
	}

	// parsing integer literal
	if unicode.IsDigit(sym) {
		number := []rune{sym}
		cursor++
		for unicode.IsDigit(src[cursor]) {
			number = append(number, src[cursor])
			cursor++
		}

		if sym = src[cursor]; !(unicode.IsSpace(sym) || sym == ')') {
			return cursor, Token{}, ErrInvalidIntegerLiteral
		}

		return cursor, Token{value: string(number), t: Int}, nil
	}

	// Check for an identifier
	if isValidChar(sym) {
		identifier := []rune{sym}
		cursor++
		for isValidChar(src[cursor]) || unicode.IsDigit(sym) {
			identifier = append(identifier, src[cursor])
			cursor++
		}

		return cursor, Token{value: string(identifier), t: Id}, nil
	}

	return cursor, Token{}, ErrInvalidSymbol
}

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

func isValidChar(sym rune) bool {
	return unicode.IsLetter(sym) || sym == '?' || sym == '!' || sym == '-' || sym == '_' || sym == '+'
}
