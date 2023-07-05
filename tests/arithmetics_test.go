package tests

import (
	"github.com/Vallghall/gopherscm/internal/core/operator"
	"github.com/Vallghall/gopherscm/internal/core/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDivision(t *testing.T) {
	a := types.NumberFrom(int64(10))
	b := types.NumberFrom(int64(2))
	expected := types.NumberFrom(int64(5))
	actual, err := a.ApplyOperation(operator.Division, b)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestMultiply(t *testing.T) {
	a := types.NumberFrom(int64(10))
	b := types.NumberFrom(int64(2))
	expected := types.NumberFrom(int64(20))
	actual, err := a.ApplyOperation(operator.Multiplication, b)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}
