package funcs_test

import (
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
	m, err := funcs.StructToMap(&s)
	require.NoError(t, err)
	expeced := map[string]interface{}{
		"age":    0,
		"belive": "",
		"id":     1,
		"name":   "name",
	}
	assert.EqualValues(t, expeced, m)
}
