package lexer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// func unwrap[T any](value T, err error) T {
// 	return value
// }

func TestLex(t *testing.T) {
	t.Run("Simple addition", func(t *testing.T) {
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
}
