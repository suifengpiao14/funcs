package funcs

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"unsafe"

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
	if IsNil(v) {
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

// StructField 结构体字段元数据定义
type StructField struct {
	structRef reflect.Value     // 结构体实例的反射值
	name      string            // 字段名
	pkgPath   string            // 包路径（用于判断是否为私有字段）
	typ       reflect.Type      // 字段类型
	Tag       reflect.StructTag // 字段标签信息，如 json:"name" 或 "-" 等
	offset    uintptr           // 偏移量（相对于结构体起始位置的字节数）
	index     []int             // 索引序列，用于通过 Type.FieldByIndex 获取嵌套的子字段
	anonymous bool              // 是否是匿名字段
}

func (f *StructField) GetName() string {
	return f.name
}

// StructFields 结构体字段列表
type StructFields []StructField

// NewStructFields 解析结构体（或结构体指针）的所有字段信息
// 参数 struc：结构体指针
// 返回值：包含所有导出字段的 StructFields 列表
func NewStructFields(struc any) (structFields StructFields) {
	if struc == nil {
		return nil
	}
	// 1. 获取结构体的反射值和类型（处理指针）
	val := reflect.ValueOf(struc)
	for val.Kind() != reflect.Ptr {
		err := errors.New("struc must be a pointer to a struct")
		panic(err)
	}
	val = val.Elem()
	typ := val.Type()
	// 校验：必须是结构体类型
	if typ.Kind() != reflect.Struct {
		err := errors.New("struc must be a pointer to a struct")
		panic(err)
	}

	// 2. 递归解析结构体字段（包括嵌套结构体）
	structFields = parseStructFields(val, typ, nil)

	return structFields
}

// parseStructFields 递归解析结构体字段
// val：当前结构体的反射值
// typ：当前结构体的反射类型
// parentIndex：父字段的索引（用于嵌套结构体）
func parseStructFields(val reflect.Value, typ reflect.Type, parentIndex []int) StructFields {
	var fields StructFields

	// 遍历所有字段
	for i := 0; i < typ.NumField(); i++ {
		fieldType := typ.Field(i)
		// 跳过未导出字段（PkgPath 非空表示私有字段）
		if fieldType.PkgPath != "" {
			continue
		}

		// 构建字段的完整索引（父索引 + 当前索引）
		fieldIndex := append(parentIndex, i)
		// 获取字段的反射值
		fieldVal := val.Field(i)

		// 3. 处理嵌套结构体（递归解析）
		elemType := fieldType.Type
		elemVal := fieldVal
		// 解引用指针类型的字段
		for elemType.Kind() == reflect.Ptr {
			elemType = elemType.Elem()
			if elemVal.Kind() == reflect.Ptr {
				if elemVal.IsNil() {
					// 如果指针字段为 nil，创建零值实例用于解析
					elemVal = reflect.New(elemType).Elem()
				} else {
					elemVal = elemVal.Elem()
				}
			}
		}

		// 如果是结构体类型且不是基础类型（如 time.Time），递归解析
		if elemType.Kind() == reflect.Struct && !isBasicStruct(elemType) {
			nestedFields := parseStructFields(elemVal, elemType, fieldIndex)
			fields = append(fields, nestedFields...)
		} else {
			if !val.CanAddr() {
				err := errors.New("struct must be addressable")
				panic(err)
			}
			// 普通字段：构建 StructField 实例
			field := StructField{
				structRef: val.Addr(), // 关联到最外层结构体实例
				name:      fieldType.Name,
				pkgPath:   fieldType.PkgPath,
				typ:       fieldType.Type,
				Tag:       fieldType.Tag,
				offset:    fieldType.Offset,
				index:     fieldIndex,
				anonymous: fieldType.Anonymous,
			}
			fields = append(fields, field)
		}
	}

	return fields
}

// isBasicStruct 判断是否是基础结构体类型（不需要递归解析）
func isBasicStruct(typ reflect.Type) bool {
	// 常见的基础结构体类型（可根据需要扩展）
	basicTypes := map[string]bool{
		"time.Time": true,
	}
	return basicTypes[typ.PkgPath()+"."+typ.Name()]
}

// GetByFieldAddress 根据字段的内存地址获取对应的字段信息
// 参数 attrAddr：字段的指针（如 &user.Addr）
// 返回值：匹配的字段信息和是否找到的标志
func (fs StructFields) GetByFieldAddress(attrAddr any) (field *StructField) {
	if attrAddr == nil || len(fs) == 0 {
		return nil
	}

	// 1. 获取字段指针的内存地址
	fieldPtr := reflect.ValueOf(attrAddr)
	if fieldPtr.Kind() != reflect.Ptr {
		err := errors.New("attrAddr must be a pointer")
		panic(err)
	}
	fieldAddr := fieldPtr.UnsafePointer()
	var structPtr unsafe.Pointer
	// 2. 遍历所有字段，匹配内存地址
	for i := range fs {
		f := &fs[i]
		// 获取结构体实例的起始地址
		if structPtr == nil {
			structPtr = f.structRef.UnsafePointer()
		}
		// 计算字段的理论地址 = 结构体起始地址 + 字段偏移量
		expectedAddr := unsafe.Add(structPtr, f.offset)

		// 地址匹配则返回该字段
		if expectedAddr == fieldAddr {
			return f
		}
	}

	err := errors.New("can not find field by address")
	panic(err)
}
