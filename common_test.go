package funcs_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/suifengpiao14/funcs"
)

func TestCreateNewInstance(t *testing.T) {

	t.Run("any", func(t *testing.T) {
		inst, isPtr := funcs.CreateNewInstance[any]()
		require.False(t, isPtr)
		require.Nil(t, inst)
	})

	t.Run("int", func(t *testing.T) {
		inst, isPtr := funcs.CreateNewInstance[int]()
		require.False(t, isPtr)
		require.NotNil(t, inst)
	})

	t.Run("*int", func(t *testing.T) {
		inst, isPtr := funcs.CreateNewInstance[*int]()
		require.True(t, isPtr)
		require.NotNil(t, inst)
	})

	t.Run("***int", func(t *testing.T) {
		inst, isPtr := funcs.CreateNewInstance[***int]()
		require.True(t, isPtr)
		require.NotNil(t, inst)
	})
}
