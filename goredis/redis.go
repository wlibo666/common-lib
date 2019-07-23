package goredis

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

var (
	redisPools = &sync.Map{}

	ERR_NOT_FOUND_CLIENT = errors.New("not found client in redis pools")
)

const (
	DEFAULT_NAME  = "default"
	MIN_POOL_SIZE = 10
)

func InitClient(addr, pwd string, db int, poolSize int, name ...string) error {
	size := poolSize
	if size < MIN_POOL_SIZE {
		size = MIN_POOL_SIZE
	}
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd,
		DB:       db,
		PoolSize: size,
	})
	_, err := client.Ping().Result()
	if err != nil {
		return fmt.Errorf("ping client:[%s] failed,err:[%s]", addr, err.Error())
	}
	redisName := DEFAULT_NAME
	if len(name) > 0 && name[0] != "" {
		redisName = name[0]
	}
	redisPools.Store(redisName, client)
	return nil
}

func GetClient(name ...string) *redis.Client {
	redisName := DEFAULT_NAME
	if len(name) > 0 && name[0] != "" {
		redisName = name[0]
	}
	v, ok := redisPools.Load(redisName)
	if !ok {
		return nil
	}
	return v.(*redis.Client)
}

func Exist(key string, name ...string) (bool, error) {
	cli := GetClient(name...)
	if cli == nil {
		return false, ERR_NOT_FOUND_CLIENT
	}
	cnt, err := cli.Exists(key).Result()
	if err != nil {
		return false, err
	}
	if cnt == 1 {
		return true, nil
	}
	return false, nil
}

func HSet(key, field, value string, name ...string) (bool, error) {
	cli := GetClient(name...)
	if cli == nil {
		return false, ERR_NOT_FOUND_CLIENT
	}
	return cli.HSet(key, field, value).Result()
}

func HGet(key, field string, name ...string) (string, error) {
	cli := GetClient(name...)
	if cli == nil {
		return "", ERR_NOT_FOUND_CLIENT
	}
	cmd := cli.HGet(key, field)
	return cmd.Result()
}

func HGetall(key string, name ...string) (map[string]string, error) {
	cli := GetClient(name...)
	if cli == nil {
		return nil, ERR_NOT_FOUND_CLIENT
	}
	cmd := cli.HGetAll(key)
	return cmd.Result()
}

func HDel(key, filed string, name ...string) (int64, error) {
	cli := GetClient(name...)
	if cli == nil {
		return 0, ERR_NOT_FOUND_CLIENT
	}
	return cli.HDel(key, filed).Result()
}

func HDels(key string, fields []string, name ...string) (int64, error) {
	cli := GetClient(name...)
	if cli == nil {
		return 0, ERR_NOT_FOUND_CLIENT
	}
	return cli.HDel(key, fields...).Result()
}

func SetString(key, value string, expirySeconds int, name ...string) error {
	cli := GetClient(name...)
	if cli == nil {
		return ERR_NOT_FOUND_CLIENT
	}
	cmd := cli.Set(key, value, time.Duration(expirySeconds)*time.Second)
	_, err := cmd.Result()
	return err
}

func Delete(key string, name ...string) error {
	cli := GetClient(name...)
	if cli == nil {
		return ERR_NOT_FOUND_CLIENT
	}
	_, err := cli.Del(key).Result()
	return err
}

func BatchDelete(keys []string, name ...string) (int64, error) {
	cli := GetClient(name...)
	if cli == nil {
		return 0, ERR_NOT_FOUND_CLIENT
	}
	return cli.Del(keys...).Result()
}

// 处理scan返回key的函数定义
type ScanProcFunc func(key string) error

func Scan(keyFlag string, count int64, scanFunc ScanProcFunc, threadNum int, name ...string) error {
	var scanPos uint64 = 0
	scanLock := sync.Mutex{}
	exitFlag := false
	wg := &sync.WaitGroup{}

	// SCAN匹配关键字
	key := keyFlag
	if key == "" {
		key = "*"
	}
	// 每次返回结果条数
	var tmpCount int64 = 1000
	if count > 0 {
		tmpCount = count
	}

	for i := 0; i < threadNum; i++ {
		wg.Add(1)
		go func(name ...string) error {
			var err error = nil
			defer wg.Done()

			redisName := DEFAULT_NAME
			if len(name) >= 1 {
				redisName = name[0]
			}
			cli := GetClient(redisName)
			if cli == nil {
				return ERR_NOT_FOUND_CLIENT
			}
			for {
				scanLock.Lock()
				if exitFlag {
					scanLock.Unlock()
					break
				}
				// 执行一次scan操作
				keys, curpos, err := cli.Scan(scanPos, key, tmpCount).Result()
				if err != nil {
					scanLock.Unlock()
					break
				}
				// 获取下一次scan位置
				if curpos > 0 {
					scanPos = curpos
				} else {
					exitFlag = true
				}
				scanLock.Unlock()
				if scanFunc != nil {
					for _, v := range keys {
						scanFunc(v)
					}
				}
			}
			return err
		}(name...)
	}
	wg.Wait()
	return nil
}

// 处理scan返回key的函数定义
type HScanProcFunc func(key, field string) error

func HScan(hashKey, fieldFlag string, count int64, scanFunc HScanProcFunc, threadNum int, name ...string) error {
	var scanPos uint64 = 0
	scanLock := sync.Mutex{}
	exitFlag := false
	wg := &sync.WaitGroup{}

	// SCAN匹配关键字
	key := fieldFlag
	if key == "" {
		key = "*"
	}
	// 每次返回结果条数
	var tmpCount int64 = 1000
	if count > 0 {
		tmpCount = count
	}

	for i := 0; i < threadNum; i++ {
		wg.Add(1)
		go func(name ...string) error {
			var err error = nil
			defer wg.Done()

			redisName := DEFAULT_NAME
			if len(name) >= 1 {
				redisName = name[0]
			}
			cli := GetClient(redisName)
			if cli == nil {
				return ERR_NOT_FOUND_CLIENT
			}
			for {
				scanLock.Lock()
				if exitFlag {
					scanLock.Unlock()
					break
				}
				// 执行一次scan操作
				keys, curpos, err := cli.HScan(hashKey, scanPos, key, tmpCount).Result()
				if err != nil {
					scanLock.Unlock()
					break
				}
				// 获取下一次scan位置
				if curpos > 0 {
					scanPos = curpos
				} else {
					exitFlag = true
				}
				scanLock.Unlock()
				if scanFunc != nil {
					for _, v := range keys {
						scanFunc(hashKey, v)
					}
				}
			}
			return err
		}(name...)
	}
	wg.Wait()
	return nil
}
