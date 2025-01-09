package funcs_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/suifengpiao14/funcs"
)

type TestStruct struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Belive string `json:"belive"`
	Mor    string `json:"-"`
}

func TestStructToMap(t *testing.T) {
	s := TestStruct{
		Id:   1,
		Name: "name",
	}
	m := funcs.StructToMap(&s)
	expeced := map[string]interface{}{
		"age":    0,
		"belive": "",
		"id":     1,
		"name":   "name",
	}
	assert.EqualValues(t, expeced, m)
}

func TestColumn(t *testing.T) {
	s := []TestStruct{
		{
			Id:   1,
			Name: "name",
		},
		{
			Id:   2,
			Name: "name",
		},
	}
	m := funcs.Column(s, func(row TestStruct) int {
		return row.Id
	})
	fmt.Println(m)
}
func TestUniqueue(t *testing.T) {
	m := funcs.Uniqueue([]int{1, 1, 3, 3, 6})
	fmt.Println(m)
}
