package utils

import (
	"bytes"
	"context"
	"crypto/md5"
	crand "crypto/rand"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// 获取指定范围内的随机数
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

// 将毫秒转换为年月日
func GetDateByMs(ms int64) string {
	t := time.Unix(ms/1000, 0)
	return fmt.Sprintf("%04d%02d%02d", t.Year(), t.Month(), t.Day())
}

// 获取环境变量
func GetEnvDef(key, defV string) string {
	v := os.Getenv(key)
	if v == "" {
		return defV
	}
	return v
}

// 设置使用最大个数CPU
func SetMaxProc() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

// 统计每个函数的使用时间，单位是纳秒
var (
	GFuncCostTime *sync.Map = &sync.Map{}
)

// 函数调用次数及耗费时间结构体
type CostTimeStore struct {
	Cnt  int64
	Cost time.Duration
}

// 函数调用结束后调用该接口,添加该函数调用时间
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

// 给指定函数添加指定调用时间.
// 此种情况是已经计算出该函数的运行时间
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

// 循环打印每个函数调用的统计信息:函数名称,调用次数,总纳秒数,平均纳秒,总毫秒数,平均毫秒
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

var (
	removePathPos int
)

func SetRemovePathPos(pos int) {
	removePathPos = pos
}

// 获取当前代码的文件名称和行数
func GetFileAndLine() string {
	_, file, line, _ := runtime.Caller(1)
	return fmt.Sprintf("%s:%d", RemovePath(file, removePathPos), line)
}

func RemovePath(name string, pos int) string {
	if removePathPos == 0 {
		return name
	}
	s := strings.Split(name, "/")
	if len(s) > pos {
		return strings.Join(s[pos:], "/")
	}
	return name
}

// 程序退出时等待秒数，等待原因是程序启动后出错退出时可能无法及时记录日志.
const (
	EXIT_WAIT_DURATION = 3
)

// 退出时指定等待时间和退出码
func ExitWait(duration, code int) {
	time.Sleep(time.Duration(duration) * time.Second)
	os.Exit(code)
}

// 退出时指定退出码
func ExitWaitDef(code int) {
	time.Sleep(time.Duration(EXIT_WAIT_DURATION) * time.Second)
	os.Exit(code)
}

// 获取指定网卡IP地址
func GetWanAddr(name string) (string, error) {
	dev, err := net.InterfaceByName(name)
	if err != nil {
		return "", err
	}
	addrs, err := dev.Addrs()
	if err != nil {
		return "", err
	}
	return strings.Split(addrs[0].String(), "/")[0], nil
}

// 根据IP和端口生成地址
func GenListenAddr(addr string, port int) string {
	if addr == "" {
		return fmt.Sprintf(":%d", port)
	}
	return fmt.Sprintf("%s:%d", addr, port)
}

// key后面添加唯一后缀
func AddNanoSufix(key string) string {
	return fmt.Sprintf("%s-%x", key, time.Now().UnixNano())
}

// 生成带超时时间的上下文
func NewContextWithTimeout(timeout int) context.Context {
	timeoutContext, _ := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	return timeoutContext
}

func GenMd5Sign(data []byte) string {
	return fmt.Sprintf("%x", md5.Sum(data))
}

func MkdirByFile(filename string) error {
	return os.MkdirAll(filepath.Dir(filename), os.ModePerm)
}

func JsonIndent(data []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, data, "", "    ")
	if err != nil {
		return []byte{}, err
	}
	return out.Bytes(), nil
}

func GenStr32() (string, error) {
	c := make([]byte, 32)
	_, e := crand.Read(c)
	if e != nil {
		return "", e
	}
	return fmt.Sprintf("%x", md5.Sum(c)), nil
}
