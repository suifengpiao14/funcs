package funcs

import (
	"errors"
	"fmt"
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
func GetCallFuncname(skip int) (funcname string) {
	var pcArr [32]uintptr // at least 1 entry needed
	var frames *runtime.Frames
	skip += 2
	n := runtime.Callers(skip, pcArr[:])
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
	funcname = fnPtr.Name()
	return funcname
}

func SplitFullFuncName(fullName string) (packageName string, funcName string) {
	lastSlash := strings.LastIndex(fullName, "/")
	packageName = fullName
	if lastSlash > -1 {
		packageName = fullName[:lastSlash+1]
		funcName = fullName[lastSlash+1:]
	}
	arr := strings.SplitN(funcName, ".", 2)
	packageName = fmt.Sprintf("%s%s", packageName, arr[0])
	funcName = arr[1]
	return packageName, funcName
}
