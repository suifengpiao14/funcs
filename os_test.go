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
	assert.Equal(t, "github.com/suifengpiao14/funcs_test.abc.func1", funcName)
}

func abc() (fn func()) {
	return func() {}
}

func TestSplitFuncName(t *testing.T) {
	packName, funcName := funcs.SplitFullFuncName("github.com/suifengpiao14/funcs_test.abc.func1")
	assert.Equal(t, "github.com/suifengpiao14/funcs_test", packName)
	assert.Equal(t, "abc.func1", funcName)
}

type User struct {
	ID string
}

func TestReflectName(t *testing.T) {
	u := new(User)
	name := funcs.GetStructName(u)
	fmt.Println(name)

}
