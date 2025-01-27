package funcs

func AppendReplace[T any](arr []T, eqFun func(first T, second T) bool, eles ...T) (newArr []T) {
	if len(arr) == 0 {
		return eles
	}
	for _, ele := range eles {
		exists := false
		for i, existsEle := range arr {
			isEqual := eqFun(existsEle, ele)
			if isEqual {
				(arr)[i] = ele
				exists = true
				break
			}
		}

		if !exists {
			(arr) = append((arr), ele)
		}
	}
	return arr
}
