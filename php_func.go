package funcs

import "strings"

func JsonEscape(value string) string { // 转义json字符串中的特殊字符，copy于mysql的escape_string函数

	dest := make([]byte, 0, 2*len(value))
	var escape byte
	for i := 0; i < len(value); i++ {
		c := value[i]
		escape = 0
		switch c {
		case 0: /* Must be escaped for 'mysql' */
			escape = '0'
		case '\n': /* Must be escaped for logs */
			escape = 'n'
		case '\r':
			escape = 'r'
		case '\t':
			escape = 't'
		case '\\':
			escape = '\\'
		case '\'':
			escape = '\''
		case '"': /* Better safe than sorry */
			escape = '"'
		case '\032': //十进制26,八进制32,十六进制1a, /* This gives problems on Win32 */
			escape = 'Z'
		}

		if escape != 0 {
			dest = append(dest, '\\', escape)
		} else {
			dest = append(dest, c)
		}
	}

	return string(dest)
}

func Addslashes(str string) string {
	var tmpRune []rune
	for _, ch := range str {
		switch ch {
		case []rune{'\\'}[0], []rune{'"'}[0], []rune{'\''}[0]:
			tmpRune = append(tmpRune, []rune{'\\'}[0])
			tmpRune = append(tmpRune, ch)
		default:
			tmpRune = append(tmpRune, ch)
		}
	}
	return string(tmpRune)
}

func Column[T any, v any](rows []T, fn func(row T) v) []v {
	var result []v
	for _, row := range rows {
		result = append(result, fn(row))
	}
	return result
}

func ColumnWithUniqueue[T any, V any](rows []T, fn func(row T) V) []V {
	var result []V
	for _, row := range rows {
		result = append(result, fn(row))
	}
	return Uniqueue(result)
}

func Uniqueue[T any](arr []T) []T {
	var result []T
	m := make(map[any]struct{})
	for _, v := range arr {
		if _, ok := m[v]; !ok {
			m[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}

func Contains[T int | int64 | string](arr []T, v T) bool {
	for _, v2 := range arr {
		if v == v2 {
			return true
		}
	}
	return false
}

func First[T any](rows []T) (first *T, exists bool) {
	if len(rows) == 0 {
		return nil, false
	}
	return &rows[0], true
}

func FirstWithDefault[T any](rows []T) (first T) {
	if len(rows) == 0 {
		return *new(T)
	}
	return rows[0]
}

func GetOne[T any](rows []T, fn func(row T) bool) (row *T, exists bool) {
	for _, r := range rows {
		if fn(r) {
			return &r, true
		}
	}
	return nil, false
}

func GetOneWithDefault[T any](rows []T, fn func(row T) bool) (row T) {
	for _, r := range rows {
		if fn(r) {
			return r
		}
	}
	return *new(T)
}

func IsEmpty[T any](rows []T) (yes bool) {
	return len(rows) == 0
}

func Filter[T any](rows []T, fn func(one T) bool) (sub []T) {
	sub = make([]T, 0)
	for _, v := range rows {
		if fn(v) {
			sub = append(sub, v)
		}
	}
	return sub
}
func Walk[T any](rows []T, fn func(one *T) (err error)) (err error) {
	for _, v := range rows {
		if err = fn(&v); err != nil {
			return err
		}
	}
	return nil
}
func Map[T any, v any](rows []T, fn func(one T) (value v)) (values []v) {
	values = make([]v, 0)
	for _, v := range rows {
		values = append(values, fn(v))
	}
	return values
}

func Strtr(str string, replace map[string]string) string {
	if len(replace) == 0 || len(str) == 0 {
		return str
	}
	for old, new := range replace {
		str = strings.ReplaceAll(str, old, new)
	}
	return str
}
