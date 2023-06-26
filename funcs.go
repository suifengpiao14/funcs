package funcs

import (
	"reflect"

	"github.com/pkg/errors"
)

func StructToMap(obj interface{}) map[string]interface{} {
	v := reflect.Indirect(reflect.ValueOf(obj))
	t := v.Type()
	if t.Kind() != reflect.Struct {
		err := errors.New("obj must be a struct or a pointer to a struct")
		panic(err)
	}
	resultMap := make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		tag := field.Tag.Get("json")
		if tag == "" || tag == "-" {
			continue // Skip fields without json tag or with "-" tag
		}
		resultMap[tag] = fieldValue.Interface()
	}
	return resultMap
}
