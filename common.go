package funcs

import "reflect"

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
	rv := reflect.ValueOf(&t).Elem() // 获取T类型的反射值

	// 构建新实例的反射值
	var newVal reflect.Value
	if rv.Kind() == reflect.Ptr {
		// 情况1：T是指针类型 → 创建指向新实例的指针
		elemType := rv.Type().Elem()     // 获取指针指向的元素类型
		newElem := reflect.New(elemType) // 创建元素的指针（*elemType）
		newVal = newElem
		isPtr = true
	} else {
		// 情况2：T是非指针类型 → 直接创建新实例
		newVal = reflect.New(rv.Type()).Elem() // 创建T类型的实例
	}
	// 将反射值转换为T类型并返回
	t = newVal.Interface().(T)
	return t, isPtr
}
