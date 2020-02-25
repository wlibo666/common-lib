package webutils

import (
	"crypto/md5"
	"fmt"
	"sort"
	"strings"
	"testing"
	"time"
)

func genSign(args map[string]string, randKey string) string {
	var params sort.StringSlice

	// 将随机key的md5值加入数组
	params = append(params, fmt.Sprintf("%x", md5.Sum([]byte(randKey))))
	// 将参数值组成数组
	for _, v := range args {
		params = append(params, v)
	}
	// 将参数以升序排序
	sort.Sort(params)
	// 参数以$连接
	paramStr := strings.Join(params, "$")
	// 计算参数md5
	return fmt.Sprintf("%x", md5.Sum([]byte(paramStr)))
}

func TestGenSignature(t *testing.T) {
	url := "/api/v1/push/register"
	host := "push-server.smartisan.com"
	method := "POST"
	// 时间戳实际使用时取当前时间，精确到秒
	now := 1573022171

	args := map[string]string{
		"method": method,
		"host":   host,
		"url":    url,
	}

	sign := genSign(args, fmt.Sprintf("%d", now))
	t.Logf("sign:%s\n", sign)
	// 实际签名 e40f218bda6ac354499904053231a930
}

func BenchmarkGenSignature(b *testing.B) {
	url := "/api/v1/push/register"
	host := "push-server.smartisan.com"
	method := "POST"
	// 时间戳实际使用时取当前时间，精确到秒

	args := map[string]string{
		"method": method,
		"host":   host,
		"url":    url,
	}
	randKey := fmt.Sprintf("%d", time.Now().Unix())
	for i := 0; i < b.N; i++ {
		genSign(args, randKey)
	}
}
