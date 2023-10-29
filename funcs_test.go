package funcs_test

import (
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
