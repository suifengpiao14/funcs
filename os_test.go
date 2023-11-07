package funcs_test

import (
	"fmt"
	"testing"

	"github.com/suifengpiao14/funcs"
)

func TestGetCallFuncname(t *testing.T) {
	funcName := funcs.GetCallFuncname()
	fmt.Println(funcName)
}

func TestGetFuncname(t *testing.T) {
	funcName := funcs.GetFuncname(abc)
	fmt.Println(funcName)
}

func abc() {}
