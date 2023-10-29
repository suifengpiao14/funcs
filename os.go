package funcs

import (
	"errors"
	"path"
	"runtime"
)

func GetRuntimeFilePath() (filepath string, err error) {
	_, filename, _, ok := runtime.Caller(1) // 跳过当前包函数
	if ok {
		return path.Dir(filename), nil
	}
	err = errors.New("not get file path by runtime.Caller(0)")
	return "", err
}
