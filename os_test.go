package funcs_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/suifengpiao14/funcs"
)

func TestGetCallFuncname(t *testing.T) {
	funcName := funcs.GetCallFuncname(0)
	fmt.Println(funcName)
}

func TestGetFuncname(t *testing.T) {
	fn := abc()
	funcName := funcs.GetFuncname(fn)
	assert.Equal(t, "funcs_test.abc.func1", funcName)
}

func abc() (fn func()) {
	return func() {}
}
