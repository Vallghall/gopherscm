package lexer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// func unwrap[t any](value t, err error) t {
// 	return value
// }

func TestLex(t *testing.T) {

	t.Run("id with a few ints", func(t *testing.T) {
		ts, err := Lex([]rune("(+ 1 2)"))
		require.NoErrorf(t, err, "expected no err, got: %v", err)
		expected := TokenStream{
			{value: "(", t: Syntax},
			{value: "+", t: Id},
			{value: "1", t: Int},
			{value: "2", t: Int},
			{value: ")", t: Syntax},
		}

		require.Equal(t, len(expected), len(ts))

		for i, tkn := range ts {
			require.Equal(t, tkn.t, expected[i].t)
			require.Equal(t, tkn.value, expected[i].value)
		}
	})

	t.Run("string tokenizing", func(t *testing.T) {
		ts, err := Lex([]rune(`(display "Hello")`))
		require.NoErrorf(t, err, "expected no err, got: %v", err)
		expected := TokenStream{
			{value: "(", t: Syntax},
			{value: "display", t: Id},
			{value: "Hello", t: String},
			{value: ")", t: Syntax},
		}

		require.Equal(t, len(expected), len(ts))

		for i, tkn := range ts {
			require.Equal(t, tkn.t, expected[i].t)
			require.Equal(t, tkn.value, expected[i].value)
		}
	})

	t.Run("nested parentheses", func(t *testing.T) {
		ts, err := Lex([]rune("(cons 1 (cons 2 (cons 3 nil)))"))
		require.NoErrorf(t, err, "expected no err, got: %v", err)

		expected := TokenStream{
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
		}

		require.Equal(t, len(expected), len(ts))

		for i, tkn := range ts {
			require.Equal(t, tkn.t, expected[i].t)
			require.Equal(t, tkn.value, expected[i].value)
		}
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
			ErrFreeClosingParantesis, err)

		_, err = Lex([]rune("((+ 1 2 3)"))
		require.ErrorIsf(
			t,
			err,
			ErrMissingClosingParenthesis,
			"expected: %v\ngot: %v",
			ErrMissingClosingParenthesis, err)
	})

	t.Run("single line comments", func(t *testing.T) {
		ts, err := Lex([]rune(`
		;; this is a comment
		(+ 1 ; 	comment
           2; 	comment again
		   3) ; also a comment
		`))
		require.NoErrorf(t, err, "expected no err, got: %v", err)

		expected := TokenStream{
			{value: "(", t: Syntax},
			{value: "+", t: Id},
			{value: "1", t: Int},
			{value: "2", t: Int},
			{value: "3", t: Int},
			{value: ")", t: Syntax},
		}

		require.Equal(t, len(expected), len(ts))

		for i, tkn := range ts {
			require.Equal(t, tkn.t, expected[i].t)
			require.Equal(t, tkn.value, expected[i].value)
		}
	})
}
