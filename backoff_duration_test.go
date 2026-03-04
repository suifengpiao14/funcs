package funcs_test

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/suifengpiao14/funcs"
)

func TestExponentialBackoff(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		err := funcs.ExponentialBackoff(funcs.DefaultBackoffDuration, funcs.NeedReDoForNetError, func() error {
			return nil
		})
		require.NoError(t, err)
	})
	t.Run("retry ok", func(t *testing.T) {
		var i = 0
		err := funcs.ExponentialBackoff(funcs.DefaultBackoffDuration, funcs.NeedReDoForNetError, func() error {
			if i > 1 {
				return nil
			}
			i++
			err := io.EOF
			return err
		})
		require.NoError(t, err)
	})

	t.Run("fail", func(t *testing.T) {
		err := funcs.ExponentialBackoff(funcs.DefaultBackoffDuration, funcs.NeedReDoForNetError, func() error {
			err := io.EOF
			return err
		})
		require.Error(t, err)
	})

}
