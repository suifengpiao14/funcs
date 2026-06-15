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
	pops := fs.Push("a")
	assert.Equal(t, []string(nil), pops)
	pops = fs.Push("b")
	assert.Equal(t, []string(nil), pops)
	pops = fs.Push("c")
	assert.Equal(t, []string(nil), pops)

	pops = fs.Push("a")
	assert.Equal(t, []string(nil), pops)

	pops = fs.Push("d")
	assert.Equal(t, []string{"a"}, pops)
}
