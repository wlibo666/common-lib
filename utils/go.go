package utils

import (
	"runtime"
	"strings"
)

func GetRoutineId() string {
	buf := make([]byte, 64)
	bufLen := runtime.Stack(buf, false)
	if bufLen < 10 {
		return ""
	}
	sps := strings.Split(string(buf), " ")
	if len(sps) != 3 {
		return ""
	}
	return sps[1]
}
