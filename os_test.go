package funcs_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/suifengpiao14/funcs"
)

func TestGetCallFuncname(t *testing.T) {
	funcName := funcs.GetCallFuncname()
	fmt.Println(funcName)
}

func TestGetFuncname(t *testing.T) {
	fn := abc
	funcName := funcs.GetFuncname(fn)
	assert.Equal(t, "abc", funcName)
}

func abc() {}
