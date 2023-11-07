package funcs_test

import (
	"fmt"
	"testing"

	"github.com/suifengpiao14/funcs"
)

func TestGetFuncname(t *testing.T) {
	funcName := funcs.GetCallFuncname()
	fmt.Println(funcName)
}
