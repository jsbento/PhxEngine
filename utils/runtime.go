package utils

import (
	"runtime"
)

func GetCurrentFuncName() string {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return "unknown"
	}
	details := runtime.FuncForPC(pc)

	return details.Name()
}

func GetCallerFuncName() string {
	pc, _, _, ok := runtime.Caller(2)
	if !ok {
		return "unknown"
	}
	details := runtime.FuncForPC(pc)

	return details.Name()
}
