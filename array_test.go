package funcs_test

import (
	"fmt"
	"testing"

	"github.com/suifengpiao14/funcs"
)

func TestAppendReplace(t *testing.T) {
	arr := []int{1, 2, 3}
	arr = funcs.AppendReplace(arr, func(first int, second int) bool {
		return first == second
	}, 4)
	fmt.Println(arr)
	arr = funcs.AppendReplace(arr, func(first int, second int) bool {
		return first == second
	}, 4, 5, 2)
	fmt.Println(arr)
}
