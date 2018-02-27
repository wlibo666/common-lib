package utils

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"time"
)

func RangeInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	interval := rand.Intn(max)
	if interval < min {
		interval += min
	}
	if interval > max {
		interval = max
	}
	return interval
}

func GetDateByMs(ms int64) string {
	t := time.Unix(ms/1000, 0)
	return fmt.Sprintf("%04d%02d%02d", t.Year(), t.Month(), t.Day())
}

func GetEnvDef(key, defV string) string {
	v := os.Getenv(key)
	if v == "" {
		return defV
	}
	return v
}

func SetMaxProc() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
