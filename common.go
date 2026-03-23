package funcs

import (
	"reflect"
)

func IsNil(v interface{}) bool {
	if v == nil {
		return true
	}
	valueOf := reflect.ValueOf(v)
	k := valueOf.Kind()
	switch k {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return valueOf.IsNil()
	default:
		return v == nil
	}
}

// CreateNewInstance 创建T类型的新实例（兼容指针/非指针类型）
func CreateNewInstance[T any]() (t T, isPtr bool) {
	var zero T
	rt := reflect.TypeOf(zero)
	// ⚠️ 关键：处理 interface / any
	if rt == nil {
		// T 是 interface{} 且为 nil
		return zero, isPtr
	}

	if rt.Kind() == reflect.Interface {
		// 无法实例化接口，返回零值
		return zero, isPtr
	}
	isPtr = rt.Kind() == reflect.Pointer
	if !isPtr {
		// T 是值类型
		t = reflect.New(rt).Elem().Interface().(T)
		return t, isPtr
	}

	// T 是指针类型
	// 处理指针类型（包括嵌套指针）：递归创建底层元素的实例，再构建指针链
	elemType := rt.Elem()
	// 一直解析到非指针的底层类型
	layer := 0
	for elemType.Kind() == reflect.Pointer {
		layer++

		elemType = elemType.Elem()
	}
	// 创建底层类型的实例，再包装回原指针层级
	newVal := reflect.New(elemType)
	// 循环条件：当前值的类型 ≠ 目标类型（未完成层级包装）
	for i := 0; i < layer; i++ {
		// 创建当前值类型的指针（*int → **int）
		newPtr := reflect.New(newVal.Type())
		// 将当前值赋值给新指针的解引用层（**int 的底层是 *int）
		newPtr.Elem().Set(newVal)
		// 迭代：当前值变为新指针
		newVal = newPtr
	}
	t = newVal.Interface().(T)
	return t, isPtr

}
