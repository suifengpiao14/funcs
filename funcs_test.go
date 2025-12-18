package funcs_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestReverse(t *testing.T) {
	m := funcs.Reverse([]int{1, 2, 3, 8})
	fmt.Println(m)
}

func TestFormatCityCode(t *testing.T) {

	t.Run("4->6", func(t *testing.T) {
		m := funcs.FormatCityCode("1234", 6)
		require.Equal(t, m, "123400")
	})
	t.Run("6->4", func(t *testing.T) {
		m := funcs.FormatCityCode("1234000", 4)
		require.Equal(t, m, "1234")
	})
	t.Run(`""->""`, func(t *testing.T) {
		m := funcs.FormatCityCode("", 4)
		require.Equal(t, m, "")
	})
}

func TestStrin(t *testing.T) {
	m := funcs.JsonString(map[string]any{})
	fmt.Println(m)
}
