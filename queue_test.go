package funcs_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/suifengpiao14/funcs"
)

func TestQueueWithLengthString(t *testing.T) {
	fs := funcs.QueueWithLength[string]{
		Length: 3,
	}
	pops, isAppend := fs.Push("a")
	assert.Equal(t, []string(nil), pops)
	assert.True(t, isAppend)
	pops, isAppend = fs.Push("b")
	assert.Equal(t, []string(nil), pops)
	assert.True(t, isAppend)
	pops, isAppend = fs.Push("c")
	assert.Equal(t, []string(nil), pops)
	assert.True(t, isAppend)
	pops, isAppend = fs.Push("a")
	assert.Equal(t, []string(nil), pops)
	assert.False(t, isAppend)

	pops, isAppend = fs.Push("d")
	assert.Equal(t, []string{"a"}, pops)
	assert.True(t, isAppend)
}
