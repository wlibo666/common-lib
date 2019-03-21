package redis

import (
	"fmt"
	"strings"
	"sync"
	"time"

	redisgo "github.com/garyburd/redigo/redis"
	"github.com/wlibo666/common-lib/atomic"
)

const (
	DEFAULT_REDIS = "default"
)

var (
	// store type: *RedisInfo
	redisInfos sync.Map
	// store type: *redisgo.Pool
	redisPools sync.Map
	// store type: * uint64
	connIndexs sync.Map
)

var (
	ERR_NOT_FOUD = fmt.Errorf("not found redis pool by name")
)

type RedisInfo struct {
	Name      string
	Addr      string
	Pwd       string
	Timeout   int
	MaxActive int
	MaxIdle   int
}

func RedisInit(addr, passwd string, timeout, active, idle int, name ...string) {
	redisName := DEFAULT_REDIS
	if len(name) >= 1 {
		redisName = name[0]
	}
	info := &RedisInfo{
		Name:      redisName,
		Addr:      addr,
		Pwd:       passwd,
		Timeout:   timeout,
		MaxActive: active,
		MaxIdle:   idle,
	}
	redisInfos.Store(redisName, info)

	redisPool := pollInit(info)
	redisPools.Store(redisName, redisPool)
}

func pollInit(info *RedisInfo) *redisgo.Pool {
	return &redisgo.Pool{
		Dial: func() (redisgo.Conn, error) {
			var err error
			addrs := strings.Split(info.Addr, ",")
			addsLen := uint64(len(addrs))
			v, _ := connIndexs.Load(info.Name)
			if v == nil {
				var index uint64
				connIndexs.Store(info.Name, &index)
				v = &index
			}

			if atomic.GetUInt64(v.(*uint64)) >= addsLen {
				atomic.ResetUInt64(v.(*uint64))
			}
			for i := atomic.GetUInt64(v.(*uint64)); i < addsLen; i++ {
				atomic.IncrUInt64(v.(*uint64))
				c, err := redisgo.Dial("tcp4", addrs[i],
					redisgo.DialConnectTimeout(time.Duration(info.Timeout)*time.Second),
				)
				if err != nil {
					continue
				}
				if info.Pwd != "" {
					_, err = c.Do("AUTH", info.Pwd)
					if err != nil {
						c.Close()
						continue
					}
					return c, nil
				}
				return c, nil
			}

			return nil, err
		},
		MaxIdle:     info.MaxIdle,
		MaxActive:   info.MaxActive,
		IdleTimeout: 60 * time.Second,
		TestOnBorrow: func(c redisgo.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func SetString(key, val string, name ...string) error {
	redisName := DEFAULT_REDIS
	if len(name) >= 1 {
		redisName = name[0]
	}
	pool, ok := redisPools.Load(redisName)
	if !ok {
		return ERR_NOT_FOUD
	}
	c := pool.(*redisgo.Pool).Get()
	defer c.Close()
	_, err := c.Do("SET", key, val)

	return err
}

func SetStringNx(key, val string, name ...string) error {
	redisName := DEFAULT_REDIS
	if len(name) >= 1 {
		redisName = name[0]
	}
	pool, ok := redisPools.Load(redisName)
	if !ok {
		return ERR_NOT_FOUD
	}
	c := pool.(*redisgo.Pool).Get()
	defer c.Close()
	_, err := c.Do("SETNX", key, val)

	return err
}

func GetString(key string, name ...string) (string, error) {
	redisName := DEFAULT_REDIS
	if len(name) >= 1 {
		redisName = name[0]
	}
	pool, ok := redisPools.Load(redisName)
	if !ok {
		return "", ERR_NOT_FOUD
	}
	c := pool.(*redisgo.Pool).Get()
	defer c.Close()

	reply, err := redisgo.String(c.Do("GET", key))
	if err != nil {
		return "", err
	}
	return reply, nil
}

func SetInt64(key string, val int64, name ...string) error {
	redisName := DEFAULT_REDIS
	if len(name) >= 1 {
		redisName = name[0]
	}
	pool, ok := redisPools.Load(redisName)
	if !ok {
		return ERR_NOT_FOUD
	}
	c := pool.(*redisgo.Pool).Get()
	defer c.Close()

	_, err := c.Do("SET", key, val)

	return err
}

func SetInt64Nx(key string, val int64, name ...string) error {
	redisName := DEFAULT_REDIS
	if len(name) >= 1 {
		redisName = name[0]
	}
	pool, ok := redisPools.Load(redisName)
	if !ok {
		return ERR_NOT_FOUD
	}
	c := pool.(*redisgo.Pool).Get()
	defer c.Close()

	_, err := c.Do("SETNX", key, val)

	return err
}

func GetInt64(key string, name ...string) (int64, error) {
	redisName := DEFAULT_REDIS
	if len(name) >= 1 {
		redisName = name[0]
	}
	pool, ok := redisPools.Load(redisName)
	if !ok {
		return 0, ERR_NOT_FOUD
	}
	c := pool.(*redisgo.Pool).Get()
	defer c.Close()

	reply, err := redisgo.Int64(c.Do("GET", key))
	if err != nil {
		return 0, err
	}
	return reply, nil
}

func HashMultiGet(name string, fields ...interface{}) ([]string, error) {
	pool, ok := redisPools.Load(name)
	if !ok {
		return []string{}, ERR_NOT_FOUD
	}
	c := pool.(*redisgo.Pool).Get()
	defer c.Close()

	return redisgo.Strings(c.Do("HMGET", fields...))
}

func HashMultiSet(name string, param ...interface{}) {
	pool, ok := redisPools.Load(name)
	if !ok {
		return
	}
	c := pool.(*redisgo.Pool).Get()
	defer c.Close()

	c.Do("HMSET", param...)
}

func HashGetAll(hashKey string, name ...string) ([]string, error) {
	redisName := DEFAULT_REDIS
	if len(name) >= 1 {
		redisName = name[0]
	}
	pool, ok := redisPools.Load(redisName)
	if !ok {
		return []string{}, ERR_NOT_FOUD
	}
	c := pool.(*redisgo.Pool).Get()
	defer c.Close()

	ret, err := c.Do("HGETALL", hashKey)
	if err != nil {
		return nil, err
	}
	if ret != nil {
		ret, err := redisgo.Strings(ret, nil)
		return ret, err
	}
	return nil, nil
}

func HashGet(hashKey string, key string, name ...string) ([]byte, error) {
	redisName := DEFAULT_REDIS
	if len(name) >= 1 {
		redisName = name[0]
	}
	pool, ok := redisPools.Load(redisName)
	if !ok {
		return []byte{}, ERR_NOT_FOUD
	}
	c := pool.(*redisgo.Pool).Get()
	defer c.Close()

	ret, err := c.Do("HGET", hashKey, key)
	if err != nil {
		return nil, err
	}
	if ret != nil {
		ret, err := redisgo.Bytes(ret, nil)
		return ret, err
	}
	return nil, nil
}

func HashSet(hashKey string, key string, val []byte, name ...string) (int, error) {
	redisName := DEFAULT_REDIS
	if len(name) >= 1 {
		redisName = name[0]
	}
	pool, ok := redisPools.Load(redisName)
	if !ok {
		return 0, ERR_NOT_FOUD
	}
	c := pool.(*redisgo.Pool).Get()
	defer c.Close()

	return redisgo.Int(c.Do("HSET", hashKey, key, val))
}

func HashDel(hashKey, field string, name ...string) (int, error) {
	redisName := DEFAULT_REDIS
	if len(name) >= 1 {
		redisName = name[0]
	}
	pool, ok := redisPools.Load(redisName)
	if !ok {
		return 0, ERR_NOT_FOUD
	}
	c := pool.(*redisgo.Pool).Get()
	defer c.Close()

	return redisgo.Int(c.Do("HDEL", hashKey, field))
}

func HashKeys(hashKey string, name ...string) ([]string, error) {
	redisName := DEFAULT_REDIS
	if len(name) >= 1 {
		redisName = name[0]
	}
	pool, ok := redisPools.Load(redisName)
	if !ok {
		return []string{}, ERR_NOT_FOUD
	}
	c := pool.(*redisgo.Pool).Get()
	defer c.Close()

	return redisgo.Strings(c.Do("HKEYS", hashKey))
}

func Incr(key string, name ...string) (int64, error) {
	redisName := DEFAULT_REDIS
	if len(name) >= 1 {
		redisName = name[0]
	}
	pool, ok := redisPools.Load(redisName)
	if !ok {
		return 0, ERR_NOT_FOUD
	}
	c := pool.(*redisgo.Pool).Get()
	defer c.Close()

	return redisgo.Int64(c.Do("INCR", key))
}

func IncrBy(key string, cnt int64, name ...string) (int64, error) {
	redisName := DEFAULT_REDIS
	if len(name) >= 1 {
		redisName = name[0]
	}
	pool, ok := redisPools.Load(redisName)
	if !ok {
		return 0, ERR_NOT_FOUD
	}
	c := pool.(*redisgo.Pool).Get()
	defer c.Close()

	return redisgo.Int64(c.Do("INCRBY", key, cnt))
}

func Expire(key string, t int64, name ...string) (int64, error) {
	redisName := DEFAULT_REDIS
	if len(name) >= 1 {
		redisName = name[0]
	}
	pool, ok := redisPools.Load(redisName)
	if !ok {
		return 0, ERR_NOT_FOUD
	}
	c := pool.(*redisgo.Pool).Get()
	defer c.Close()

	return redisgo.Int64(c.Do("EXPIRE", key, t))
}

func Ttl(key string, name ...string) (int64, error) {
	redisName := DEFAULT_REDIS
	if len(name) >= 1 {
		redisName = name[0]
	}
	pool, ok := redisPools.Load(redisName)
	if !ok {
		return 0, ERR_NOT_FOUD
	}
	c := pool.(*redisgo.Pool).Get()
	defer c.Close()

	return redisgo.Int64(c.Do("TTL", key))
}

func Delete(key string, name ...string) (int64, error) {
	redisName := DEFAULT_REDIS
	if len(name) >= 1 {
		redisName = name[0]
	}
	pool, ok := redisPools.Load(redisName)
	if !ok {
		return 0, ERR_NOT_FOUD
	}
	c := pool.(*redisgo.Pool).Get()
	defer c.Close()

	return redisgo.Int64(c.Do("DEL", key))
}

func Exist(name string, key ...interface{}) (int64, error) {
	pool, ok := redisPools.Load(name)
	if !ok {
		return 0, ERR_NOT_FOUD
	}
	c := pool.(*redisgo.Pool).Get()
	defer c.Close()

	return redisgo.Int64(c.Do("EXISTS", key...))
}

// 处理scan返回key的函数定义
type ScanProcFunc func(key string) error

func Scan(keyFlag string, count int, scanFunc ScanProcFunc, threadNum int, name ...string) error {
	var scanPos int64 = 0
	scanLock := sync.Mutex{}
	exitFlag := false
	wg := &sync.WaitGroup{}

	// SCAN匹配关键字
	key := keyFlag
	if key == "" {
		key = "*"
	}
	// 每次返回结果条数
	tmpCount := 10
	if count > 0 {
		tmpCount = count
	}

	for i := 0; i < threadNum; i++ {
		index := i
		wg.Add(1)
		go func(index int, name ...string) error {
			defer wg.Done()
			redisName := DEFAULT_REDIS
			if len(name) >= 1 {
				redisName = name[0]
			}
			pool, ok := redisPools.Load(redisName)
			if !ok {
				return ERR_NOT_FOUD
			}
			c := pool.(*redisgo.Pool).Get()
			defer c.Close()

			for {
				scanLock.Lock()
				if exitFlag {
					scanLock.Unlock()
					return nil
				}
				// 执行一次scan操作
				tmpRes, err := c.Do("SCAN", scanPos, "MATCH", key, "COUNT", tmpCount)
				if err != nil {
					scanLock.Unlock()
					return fmt.Errorf("scan thread:%d scan failed,err:%s", index, err.Error())
				}
				res := tmpRes.([]interface{})
				// 返回结果数组长度不为2,则格式错误
				if len(res) != 2 {
					scanLock.Unlock()
					return fmt.Errorf("scan thread:%d scan res len is:%d,not 2", index, len(res))
				}
				// 获取下一次scan位置
				var tmpNum int64 = 0
				fmt.Sscanf(string(res[0].([]byte)), "%d", &tmpNum)
				if tmpNum > 0 {
					scanPos = tmpNum
				} else {
					exitFlag = true
				}
				scanLock.Unlock()
				if scanFunc != nil {
					for _, v := range res[1].([]interface{}) {
						go func(key string) {
							scanFunc(key)
						}(string(v.([]byte)))
					}
				}
			}
			return nil
		}(index, name...)
	}
	wg.Wait()
	return nil
}

// 处理scan返回key的函数定义
type HScanProcFunc func(key, field string) error

func HScan(hashKey, fieldFlag string, count int, scanFunc HScanProcFunc, threadNum int, name ...string) error {
	var scanPos int64 = 0
	scanLock := sync.Mutex{}
	exitFlag := false
	wg := &sync.WaitGroup{}

	// SCAN匹配关键字
	key := fieldFlag
	if key == "" {
		key = "*"
	}
	// 每次返回结果条数
	tmpCount := 10
	if count > 0 {
		tmpCount = count
	}

	for i := 0; i < threadNum; i++ {
		in := i
		wg.Add(1)
		go func(index int, name ...string) error {
			defer wg.Done()
			redisName := DEFAULT_REDIS
			if len(name) >= 1 {
				redisName = name[0]
			}
			pool, ok := redisPools.Load(redisName)
			if !ok {
				return ERR_NOT_FOUD
			}
			c := pool.(*redisgo.Pool).Get()
			defer c.Close()

			for {
				scanLock.Lock()
				if exitFlag {
					scanLock.Unlock()
					return nil
				}
				// 执行一次scan操作
				tmpRes, err := c.Do("HSCAN", hashKey, scanPos, "MATCH", key, "COUNT", tmpCount)
				if err != nil {
					scanLock.Unlock()
					return fmt.Errorf("hscan thread:%d scan failed,err:%s", index, err.Error())
				}
				res := tmpRes.([]interface{})
				// 返回结果数组长度不为2,则格式错误
				if len(res) != 2 {
					scanLock.Unlock()
					return fmt.Errorf("hscan thread:%d scan res len is:%d,not 2", index, len(res))
				}
				// 获取下一次scan位置
				var tmpNum int64 = 0
				fmt.Sscanf(string(res[0].([]byte)), "%d", &tmpNum)
				if tmpNum > 0 {
					scanPos = tmpNum
				} else {
					exitFlag = true
				}
				scanLock.Unlock()
				if scanFunc != nil {
					v2 := res[1].([]interface{})
					resLen := len(v2)
					if resLen > 0 {
						for i := 0; i < resLen; i += 2 {
							go func(key, field string) {
								scanFunc(key, field)
							}(hashKey, string(v2[i].([]byte)))
						}
					}
				}
			}
			return nil
		}(in, name...)
	}
	wg.Wait()
	return nil
}
