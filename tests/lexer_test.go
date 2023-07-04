package tests

import (
	"testing"

	"github.com/Vallghall/gopherscm/internal/data"
	"github.com/Vallghall/gopherscm/internal/lexer"

	"github.com/stretchr/testify/require"
)

// func unwrap[t any](value t, err error) t {
// 	return value
// }

func TestLex(t *testing.T) {

	t.Run("id with a few ints", func(t *testing.T) {
		ts, err := lexer.Lex([]rune("(+ 1 2)"))
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

	t.Run("negative integers lexing", func(t *testing.T) {
		ts, err := lexer.Lex([]rune("(+ -1 2 -3)"))
		require.NoErrorf(t, err, "expected no err, got: %v", err)
		expected := data.TokenStream{
			data.NewToken("(", data.Syntax),
			data.NewToken("+", data.Id),
			data.NewToken("-1", data.Int),
			data.NewToken("2", data.Int),
			data.NewToken("-3", data.Int),
			data.NewToken(")", data.Syntax),
		}
		require.Equal(t, len(expected), len(ts))

		for i, tkn := range ts {
			require.Equal(t, tkn.Type(), expected[i].Type())
			require.Equal(t, tkn.Value(), expected[i].Value())
		}
	})

	t.Run("negative NaN", func(t *testing.T) {
		_, err := lexer.Lex([]rune(`(+ 1 -"2")`))
		require.ErrorIsf(
			t,
			err,
			lexer.ErrNaN,
			"expected: %v,\ngot: %v",
			lexer.ErrNaN, err)

		_, err = lexer.Lex([]rune(`(+ 1 -foo)`))
		require.ErrorIsf(
			t,
			err,
			lexer.ErrNaN,
			"expected: %v,\ngot: %v",
			lexer.ErrNaN, err)
	})

	t.Run("floats", func(t *testing.T) {
		ts, err := lexer.Lex([]rune("(+ -1.342 20.576 -3.0)"))
		require.NoErrorf(t, err, "expected no err, got: %v", err)
		expected := data.TokenStream{
			data.NewToken("(", data.Syntax),
			data.NewToken("+", data.Id),
			data.NewToken("-1.342", data.Float),
			data.NewToken("20.576", data.Float),
			data.NewToken("-3.0", data.Float),
			data.NewToken(")", data.Syntax),
		}
		require.Equal(t, len(expected), len(ts))

		for i, tkn := range ts {
			require.Equal(t, tkn.Type(), expected[i].Type())
			require.Equal(t, tkn.Value(), expected[i].Value())
		}
	})

	t.Run("string tokenizing", func(t *testing.T) {
		ts, err := lexer.Lex([]rune(`(display "Hello")`))
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

	t.Run("string unexpected line break", func(t *testing.T) {
		_, err := lexer.Lex([]rune(`(display "Hel
		lo")`))
		require.ErrorIsf(
			t,
			err,
			lexer.ErrUnexpectedLineBreak,
			"expected: %v,\ngot: %v",
			lexer.ErrUnexpectedLineBreak, err)

	})

	t.Run("nested parentheses", func(t *testing.T) {
		ts, err := lexer.Lex([]rune("(cons 1 (cons 2 (cons 3 nil)))"))
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
		_, err := lexer.Lex([]rune(`(+ 1 2 3`))
		require.ErrorIsf(
			t,
			err,
			lexer.ErrMissingClosingParenthesis,
			"expected: %v,\ngot: %v",
			lexer.ErrMissingClosingParenthesis, err)

		_, err = lexer.Lex([]rune("(+ 1 2 3))"))
		require.ErrorIsf(
			t,
			err,
			lexer.ErrFreeClosingParantesis,
			"expected: %v\ngot: %v",
			lexer.ErrFreeClosingParantesis, err)

		_, err = lexer.Lex([]rune("((+ 1 2 3)"))
		require.ErrorIsf(
			t,
			err,
			lexer.ErrMissingClosingParenthesis,
			"expected: %v\ngot: %v",
			lexer.ErrMissingClosingParenthesis, err)
	})

	t.Run("single line comments", func(t *testing.T) {
		ts, err := lexer.Lex([]rune(`
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
