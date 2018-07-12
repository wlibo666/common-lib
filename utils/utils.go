package utils

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"sync"
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

var (
	GFuncCostTime *sync.Map = &sync.Map{}
)

type CostTimeStore struct {
	Cnt  int64
	Cost time.Duration
}

func CostTime(funcName string, start time.Time) {
	terminal := time.Since(start)
	v, ok := GFuncCostTime.Load(funcName)
	if ok {
		v.(*CostTimeStore).Cnt += 1
		v.(*CostTimeStore).Cost += terminal
		return
	}
	tmp := &CostTimeStore{
		Cnt:  1,
		Cost: terminal,
	}
	GFuncCostTime.Store(funcName, tmp)
}

func AddCostTime(funcName string, cost time.Duration) {
	v, ok := GFuncCostTime.Load(funcName)
	if ok {
		v.(*CostTimeStore).Cnt += 1
		v.(*CostTimeStore).Cost += cost
		return
	}
	tmp := &CostTimeStore{
		Cnt:  1,
		Cost: cost,
	}
	GFuncCostTime.Store(funcName, tmp)
}

func LoopPrintCostTime(interval ...int) {
	go func(interval ...int) {
		inter := 5
		if len(interval) > 0 {
			inter = interval[0]
		}
		for {
			time.Sleep(time.Duration(inter) * time.Second)
			GFuncCostTime.Range(func(key, value interface{}) bool {
				store := value.(*CostTimeStore)
				fmt.Fprintf(os.Stdout, "func:%-30s | cnt:%-8d | cost: %-16d ns | avg:%-16d ns | %-8d ms | avg:%-8d ms\n",
					key, store.Cnt, store.Cost, store.Cost/time.Duration(store.Cnt), store.Cost/time.Millisecond,
					(store.Cost/time.Millisecond)/time.Duration(store.Cnt))
				GFuncCostTime.Delete(key)
				return true
			})
		}
	}(interval...)
}
