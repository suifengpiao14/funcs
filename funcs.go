package funcs

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

// FormatCityCode 城市id 经常有4位、6位格式，这里做一个格式化处理，保证长度一致
func FormatCityCode(cityCode string, size int) string {
	if cityCode == "" {
		return ""
	}
	s := fmt.Sprintf(`%s%s`, cityCode, strings.Repeat("0", size))[:size]
	return s
}

func JsonString(v any) string {
	if v == nil {
		return ""
	}
	b, err := json.Marshal(v)
	if err != nil {
		return err.Error()
	}
	s := string(b)
	return s
}

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

// Struct2MapString 结构体转json
// 使用场景 resty curl 构造请求参数时常用
func Struct2MapString(i interface{}) (out map[string]string, err error) {

	var myMap map[string]interface{}
	data, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &myMap)
	if err != nil {
		return nil, err
	}
	out = make(map[string]string)
	for k, v := range myMap {
		out[k] = cast.ToString(v)
	}
	return out, nil
}

func Struct2JsonMap(rows any) map[string]any {
	data, err := json.Marshal(rows)
	if err != nil {
		panic(err)
	}
	m := make(map[string]any)
	err = json.Unmarshal(data, &m)
	if err != nil {
		panic(err)
	}
	return m
}

// 拷贝 sync.oncefunc.go 低版本go 不支持 go 1.21 版本才有，直接复制
// OnceValue returns a function that invokes f only once and returns the value
// returned by f. The returned function may be called concurrently.
//
// If f panics, the returned function will panic with the same value on every call.
func OnceValue[T any](f func() T) func() T {
	var (
		once   sync.Once
		valid  bool
		p      any
		result T
	)
	g := func() {
		defer func() {
			p = recover()
			if !valid {
				panic(p)
			}
		}()
		result = f()
		valid = true
	}
	return func() T {
		once.Do(g)
		if !valid {
			panic(p)
		}
		return result
	}
}
