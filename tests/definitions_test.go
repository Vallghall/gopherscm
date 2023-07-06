package tests

import (
	"fmt"
	"github.com/Vallghall/gopherscm/internal/lexer"
	"github.com/Vallghall/gopherscm/internal/parser"
	"github.com/stretchr/testify/require"
	"math"
	"testing"
)

func TestDefine(t *testing.T) {

	t.Run("nested define", func(t *testing.T) {
		code := fmt.Sprintf(`
(define (area r)
	(define pi %v)
	(* pi r r))
(area 5)`, math.Pi)
		
		ts, err := lexer.Lex([]rune(code))
		require.NoError(t, err)

		result, err := parser.Parse(ts).Eval()
		require.NoError(t, err)

		expected := math.Pi * 25.0

		require.Equal(t, expected, result.Value())
	})
}
