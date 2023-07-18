package lexer

import (
	"errors"
	"unicode"

	"github.com/Vallghall/gopherscm/internal/errscm"

	"github.com/Vallghall/gopherscm/internal/data"
)

// Lex transforms input slice of symbols into slice of valid Scheme tokens
func Lex(src []rune) (data.TokenStream, error) {
	ts := make(data.TokenStream, 0)
	m := data.NewMeta()
	inputLength := len(src)
	cursor := 0

	var token *data.Token
	var err error
	parenCount := 0

	// cursor is updated in Tokenize and skip* functions only
	for cursor < inputLength {
		old := cursor
		cursor = skipSingleLineComment(cursor, src, m)
		cursor = skipWhiteSpaces(cursor, src, m)
		if old != cursor {
			continue // recheck whitespaces and comments after their deletion
		}

		cursor, token, err = Tokenize(cursor, src, m)
		if err != nil {
			if errors.Is(err, errscm.ErrEndOfInput) {
				break
			}
			return nil, err
		}

		ts = append(ts, token)
		if token.Type() == data.Syntax {
			if token.Value() == "(" {
				parenCount++
			} else {
				parenCount--
			}
		}

		if parenCount < 0 {
			return nil, errscm.ErrFreeClosingParenthesis
		}
	}

	if parenCount > 0 {
		return nil, errscm.ErrMissingClosingParenthesis
	}

	return ts, nil
}

// Tokenize - Extracts token from rune sequence
func Tokenize(cursor int, src []rune, m *data.Meta) (int, *data.Token, error) {
	if cursor >= len(src) {
		return cursor, nil, errscm.ErrEndOfInput
	}

	sym := src[cursor]

	// \' is a special sugar for `quote`` builtin func
	if sym == '\'' {
		t := data.TokenFromMeta(m)
		m.Inc()
		return cursor + 1, t.Set(data.Quote, sym), nil
	}

	// '(' and ')' are the only syntax tokens
	if sym == '(' || sym == ')' {
		t := data.TokenFromMeta(m)
		m.Inc()
		return cursor + 1, t.Set(data.Syntax, sym), nil
	}

	// parsing string literal like "foo"
	if sym == '"' {
		return extractString(cursor, src, m)
	}

	// parsing integer literal
	// TODO: add floating point numbers lexing
	if unicode.IsDigit(sym) || sym == '-' {
		return extractNumber(cursor, src, m)
	}

	// Check for an identifier
	if isValidChar(sym) {
		return extractIdentifier(cursor, src, m)
	}

	return cursor, nil, errscm.ErrInvalidSymbol
}

// skipSingleLineComment
func skipSingleLineComment(cursor int, src []rune, m *data.Meta) int {
	inputLength := len(src)
	if src[cursor] == ';' {
		for src[cursor] != '\n' {
			m.Inc()
			cursor++
			if cursor >= inputLength {
				return cursor
			}
		}
		m.NewLine()
		cursor++
	}

	return cursor
}

// skipWhitespaces - helper func for omitting whitespaces
func skipWhiteSpaces(cursor int, src []rune, m *data.Meta) int {
	inputLength := len(src)
	for unicode.IsSpace(src[cursor]) {
		cursor++
		if cursor >= inputLength {
			return cursor
		}

		m.IncNL(src[cursor])
	}

	return cursor
}

// extractString - helper func for extracting String token
// TODO: add support for escape sequences
func extractString(cursor int, src []rune, m *data.Meta) (int, *data.Token, error) {
	t := data.TokenFromMeta(m)
	cursor++ // move forward from quote
	if cursor >= len(src) {
		return cursor, nil, errscm.ErrEndOfInput
	}
	m.Inc()

	str := make([]rune, 0)
	for sym := src[cursor]; sym != '"'; sym = src[cursor] {
		if sym == '\n' {
			return cursor, nil, errscm.ErrUnexpectedLineBreak
		}

		str = append(str, sym)
		cursor++
		if cursor >= len(src) {
			return cursor, nil, errscm.ErrMissingMatchingDoubleQuotes
		}

		m.Inc()
	}

	m.Inc()
	cursor++
	return cursor, t.Set(data.String, str...), nil
}

// extractNumber - helper func for extracting an Int token
// TODO: add floating point scientific notation support
func extractNumber(cursor int, src []rune, m *data.Meta) (int, *data.Token, error) {
	t := data.TokenFromMeta(m)
	isFloat := false

	number := []rune{src[cursor]}
	cursor++
	if cursor >= len(src) {
		return cursor, nil, errscm.ErrEndOfInput
	}

	// check situations like -foo or -"foo"
	if number[0] == '-' && !unicode.IsDigit(src[cursor]) && src[cursor-2] != '(' {
		return cursor - 1, nil, errscm.ErrNaN
	}

	m.Inc()

	for unicode.IsDigit(src[cursor]) || src[cursor] == '.' {
		number = append(number, src[cursor])
		if src[cursor] == '.' {
			if isFloat {
				return cursor, nil, errscm.ErrUnexpectedDotSymbol
			}
			isFloat = true
		}

		cursor++
		if cursor >= len(src) {
			return cursor, nil, errscm.ErrEndOfInput
		}

		m.Inc()
	}

	cursor = skipSingleLineComment(cursor, src, m)
	if cursor >= len(src) {
		return cursor, nil, errscm.ErrEndOfInput
	}

	if sym := src[cursor]; !(unicode.IsSpace(sym) || sym == ')') {
		return cursor, nil, errscm.ErrInvalidNumericLiteral
	}

	if isFloat {
		return cursor, t.Set(data.Float, number...), nil
	}
	return cursor, t.Set(data.Int, number...), nil
}

// extractIdentifier - helper func for lexing identifiers
func extractIdentifier(cursor int, src []rune, m *data.Meta) (int, *data.Token, error) {
	t := data.TokenFromMeta(m)
	id := []rune{src[cursor]}
	cursor++
	if cursor >= len(src) {
		return cursor, nil, errscm.ErrEndOfInput
	}
	m.IncNL(src[cursor])

	for isValidChar(src[cursor]) || unicode.IsDigit(src[cursor]) {
		id = append(id, src[cursor])

		cursor++
		if cursor >= len(src) {
			return cursor, nil, errscm.ErrEndOfInput
		}
		m.IncNL(src[cursor])
	}

	return cursor, t.Set(data.Id, id...), nil
}

// isValidChar - predicate for checking a valid identifier's symbol
func isValidChar(sym rune) bool {
	return unicode.IsLetter(sym) ||
		sym == '?' || sym == '!' ||
		sym == '-' || sym == '_' ||
		sym == '+' || sym == '*' ||
		sym == '/'
}
