package lexer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// func unwrap[T any](value T, err error) T {
// 	return value
// }

func TestLex(t *testing.T) {

	t.Run("id with a few ints", func(t *testing.T) {
		ts, err := Lex([]rune("(+ 1 2)"))
		require.NoErrorf(t, err, "expected no err, got: %v", err)

		require.Equal(t, []Token{
			{value: "(", t: Syntax},
			{value: "+", t: Id},
			{value: "1", t: Int},
			{value: "2", t: Int},
			{value: ")", t: Syntax},
		}, ts)
	})

	t.Run("string tokenizing", func(t *testing.T) {
		ts, err := Lex([]rune(`(display "Hello")`))
		require.NoErrorf(t, err, "expected no err, got: %v", err)

		require.Equal(t, []Token{
			{value: "(", t: Syntax},
			{value: "display", t: Id},
			{value: "Hello", t: String},
			{value: ")", t: Syntax},
		}, ts)
	})

	t.Run("nested parentheses", func(t *testing.T) {
		ts, err := Lex([]rune("(cons 1 (cons 2 (cons 3 nil)))"))
		require.NoErrorf(t, err, "expected no err, got: %v", err)

		require.Equal(t, []Token{
			{value: "(", t: Syntax},
			{value: "cons", t: Id},
			{value: "1", t: Int},
			{value: "(", t: Syntax},
			{value: "cons", t: Id},
			{value: "2", t: Int},
			{value: "(", t: Syntax},
			{value: "cons", t: Id},
			{value: "3", t: Int},
			{value: "nil", t: Id},
			{value: ")", t: Syntax},
			{value: ")", t: Syntax},
			{value: ")", t: Syntax},
		}, ts)
	})

	t.Run("parenthesis mismatch", func(t *testing.T) {
		_, err := Lex([]rune(`(+ 1 2 3`))
		require.ErrorIsf(
			t,
			err,
			ErrMissingClosingParenthesis,
			"expected: %v,\ngot: %v",
			ErrMissingClosingParenthesis, err)

		_, err = Lex([]rune("(+ 1 2 3))"))
		require.ErrorIsf(
			t,
			err,
			ErrFreeClosingParantesis,
			"expected: %v\ngot: %v",
			ErrMissingClosingParenthesis, err)
	})
}
