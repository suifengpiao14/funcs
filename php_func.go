package funcs

import "strings"

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

func Uniqueue[T int | int64 | string](arr []T) []T {
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

func Strtr(str string, replace map[string]string) string {
	if len(replace) == 0 || len(str) == 0 {
		return str
	}
	for old, new := range replace {
		str = strings.ReplaceAll(str, old, new)
	}
	return str
}
