package funcs

import (
	"errors"
	"path"
	"reflect"
	"runtime"
	"strings"
)

func GetRuntimeFilePath() (filepath string, err error) {
	_, filename, _, ok := runtime.Caller(1) // 跳过当前包函数
	if !ok {
		err = errors.New("not get file path by runtime.Caller(0)")
		return "", err

	}
	return path.Dir(filename), nil
}

// GetCallFuncname 获取调用当前函数的名称
func GetCallFuncname() (funcname string) {
	var pcArr [32]uintptr // at least 1 entry needed
	var frames *runtime.Frames
	n := runtime.Callers(2, pcArr[:])
	frames = runtime.CallersFrames(pcArr[:n])
	frame, _ := frames.Next()
	funcname = frame.Function
	return funcname
}

func GetFuncname(fn any) (funcname string) {
	fnPtr := runtime.FuncForPC(reflect.ValueOf(fn).Pointer())
	if fnPtr == nil {
		return ""
	}
	fullName := fnPtr.Name()
	lastSlash := strings.LastIndex(fullName, "/")
	funcname = fullName
	if lastSlash > -1 {
		funcname = fullName[lastSlash+1:]
	}
	return funcname

}
