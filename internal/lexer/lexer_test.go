package lexer

import (
	"github.com/Vallghall/gopherscm/internal/data"
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
		expected := data.TokenStream{
			data.NewToken("(", data.Syntax),
			data.NewToken("+", data.Id),
			data.NewToken("1", data.Int),
			data.NewToken("2", data.Int),
			data.NewToken(")", data.Syntax),
		}
		require.Equal(t, len(expected), len(ts))

		for i, tkn := range ts {
			require.Equal(t, tkn.Type(), expected[i].Type())
			require.Equal(t, tkn.Value(), expected[i].Value())
		}
	})

	t.Run("string tokenizing", func(t *testing.T) {
		ts, err := Lex([]rune(`(display "Hello")`))
		require.NoErrorf(t, err, "expected no err, got: %v", err)
		expected := data.TokenStream{
			data.NewToken("(", data.Syntax),
			data.NewToken("display", data.Id),
			data.NewToken("Hello", data.String),
			data.NewToken(")", data.Syntax),
		}

		require.Equal(t, len(expected), len(ts))

		for i, tkn := range ts {
			require.Equal(t, tkn.Type(), expected[i].Type())
			require.Equal(t, tkn.Value(), expected[i].Value())
		}
	})

	t.Run("nested parentheses", func(t *testing.T) {
		ts, err := Lex([]rune("(cons 1 (cons 2 (cons 3 nil)))"))
		require.NoErrorf(t, err, "expected no err, got: %v", err)

		expected := data.TokenStream{
			data.NewToken("(", data.Syntax),
			data.NewToken("cons", data.Id),
			data.NewToken("1", data.Int),
			data.NewToken("(", data.Syntax),
			data.NewToken("cons", data.Id),
			data.NewToken("2", data.Int),
			data.NewToken("(", data.Syntax),
			data.NewToken("cons", data.Id),
			data.NewToken("3", data.Int),
			data.NewToken("nil", data.Id),
			data.NewToken(")", data.Syntax),
			data.NewToken(")", data.Syntax),
			data.NewToken(")", data.Syntax),
		}

		require.Equal(t, len(expected), len(ts))

		for i, tkn := range ts {
			require.Equal(t, tkn.Type(), expected[i].Type())
			require.Equal(t, tkn.Value(), expected[i].Value())
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

		expected := data.TokenStream{
			data.NewToken("(", data.Syntax),
			data.NewToken("+", data.Id),
			data.NewToken("1", data.Int),
			data.NewToken("2", data.Int),
			data.NewToken("3", data.Int),
			data.NewToken(")", data.Syntax),
		}

		require.Equal(t, len(expected), len(ts))

		for i, tkn := range ts {
			require.Equal(t, tkn.Type(), expected[i].Type())
			require.Equal(t, tkn.Value(), expected[i].Value())
		}
	})
}
