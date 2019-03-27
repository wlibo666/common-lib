package rediscluster

import (
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

type RedisCluster struct {
	client *redis.ClusterClient
}

var (
	DefaultCluster *RedisCluster
	defPoolSize    int = 500
)

func ClusterInit(addrs, pwd string, poolsize int) {
	defPoolSize = poolsize
	DefaultCluster = &RedisCluster{
		client: newClusterClient(addrs, pwd, poolsize),
	}
}

func NewRedisCluster(addrs, pwd string, poolsize int) *RedisCluster {
	cluster := &RedisCluster{
		client: newClusterClient(addrs, pwd, poolsize),
	}
	return cluster
}

func newClusterClient(addrs, pwd string, poolsize int) *redis.ClusterClient {
	ops := &redis.ClusterOptions{
		Addrs:      strings.Split(addrs, ","),
		ReadOnly:   true,
		MaxRetries: 5,
		Password:   pwd,

		DialTimeout:  time.Duration(3) * time.Second,
		ReadTimeout:  time.Duration(60) * time.Second,
		WriteTimeout: time.Duration(60) * time.Second,

		PoolSize:    poolsize,
		PoolTimeout: time.Duration(30) * time.Second,
		IdleTimeout: time.Duration(30) * time.Second,
	}
	return redis.NewClusterClient(ops)
}

func boolCmdRes(cmd *redis.BoolCmd) (bool, error) {
	return cmd.Val(), cmd.Err()
}

func stringCmdRes(cmd *redis.StringCmd) (string, error) {
	return cmd.Val(), cmd.Err()
}

func stringCmdInt64Res(cmd *redis.StringCmd) (int64, error) {
	return cmd.Int64()
}

func intCmdRes(cmd *redis.IntCmd) (int64, error) {
	return cmd.Val(), cmd.Err()
}

func duraCmdRes(cmd *redis.DurationCmd) (int64, error) {
	return int64(cmd.Val()) / int64(time.Second), cmd.Err()
}

func strSliceCmdRes(cmd *redis.StringSliceCmd) ([]string, error) {
	return cmd.Result()
}

func statusCmdRes(cmd *redis.StatusCmd) (string, error) {
	return cmd.Val(), cmd.Err()
}

func (c *RedisCluster) SetNx(key string, val interface{}, expire int64) (bool, error) {
	return boolCmdRes(c.client.SetNX(key, val, time.Duration(expire)*time.Second))
}

func SetNx(key string, val interface{}, expire int64) (bool, error) {
	return boolCmdRes(DefaultCluster.client.SetNX(key, val, time.Duration(expire)*time.Second))
}

func (c *RedisCluster) Set(key string, val interface{}, expire int64) (string, error) {
	return statusCmdRes(c.client.Set(key, val, time.Duration(expire)*time.Second))
}

func Set(key string, val interface{}, expire int64) (string, error) {
	return statusCmdRes(DefaultCluster.client.Set(key, val, time.Duration(expire)*time.Second))
}

func (c *RedisCluster) Expire(key string, expire int64) (bool, error) {
	return boolCmdRes(c.client.Expire(key, time.Duration(expire)*time.Second))
}

func Expire(key string, expire int64) (bool, error) {
	return boolCmdRes(DefaultCluster.client.Expire(key, time.Duration(expire)*time.Second))
}

func GetStr(key string) (string, error) {
	return stringCmdRes(DefaultCluster.client.Get(key))
}

func (c *RedisCluster) GetStr(key string) (string, error) {
	return stringCmdRes(c.client.Get(key))
}

func GetInt64(key string) (int64, error) {
	return stringCmdInt64Res(DefaultCluster.client.Get(key))
}

func (c *RedisCluster) GetInt64(key string) (int64, error) {
	return stringCmdInt64Res(c.client.Get(key))
}

func Del(key ...string) (int64, error) {
	return intCmdRes(DefaultCluster.client.Del(key...))
}

func (c *RedisCluster) Del(key ...string) (int64, error) {
	return intCmdRes(c.client.Del(key...))
}

func TTL(key string) (int64, error) {
	return duraCmdRes(DefaultCluster.client.TTL(key))
}

func (c *RedisCluster) TTL(key string) (int64, error) {
	return duraCmdRes(c.client.TTL(key))
}

func IncrBy(key string, value int64) (int64, error) {
	return intCmdRes(DefaultCluster.client.IncrBy(key, value))
}

func (c *RedisCluster) IncrBy(key string, value int64) (int64, error) {
	return intCmdRes(c.client.IncrBy(key, value))
}

func HincrBy(key, field string, value int64) (int64, error) {
	return intCmdRes(DefaultCluster.client.HIncrBy(key, field, value))
}

func (c *RedisCluster) HincrBy(key, field string, value int64) (int64, error) {
	return intCmdRes(c.client.HIncrBy(key, field, value))
}

func ZAdd(key string, score float64, members ...interface{}) (int64, error) {
	var sz []redis.Z
	for _, m := range members {
		z := redis.Z{
			Score:  score,
			Member: m,
		}
		sz = append(sz, z)
	}
	return intCmdRes(DefaultCluster.client.ZAdd(key, sz...))
}

func (c *RedisCluster) ZAdd(key string, score float64, members ...interface{}) (int64, error) {
	var sz []redis.Z
	for _, m := range members {
		z := redis.Z{
			Score:  score,
			Member: m,
		}
		sz = append(sz, z)
	}
	return intCmdRes(c.client.ZAdd(key, sz...))
}

func ZRange(key string, start, stop int64) ([]string, error) {
	return strSliceCmdRes(DefaultCluster.client.ZRange(key, start, stop))
}

func (c *RedisCluster) ZRange(key string, start, stop int64) ([]string, error) {
	return strSliceCmdRes(c.client.ZRange(key, start, stop))
}

func ZRemRangeByRank(key string, start, stop int64) (int64, error) {
	return intCmdRes(DefaultCluster.client.ZRemRangeByRank(key, start, stop))
}

func (c *RedisCluster) ZRemRangeByRank(key string, start, stop int64) (int64, error) {
	return intCmdRes(c.client.ZRemRangeByRank(key, start, stop))
}

func SetSortSet(key string, values interface{}, score float64) {
	DefaultCluster.client.ZAdd(key,
		redis.Z{
			Score:  score,
			Member: values})
	DefaultCluster.client.ZRemRangeByRank(key, 1, -2)
}

func GetSortSet(key string) *redis.ZSliceCmd {
	res := DefaultCluster.client.ZRangeWithScores(key, 0, 10)
	return res
}

// 处理scan返回key的函数定义
type ScanProcFunc func(key string) error

type ScanCmd struct {
	Match     string
	Count     int64
	ScanFunc  ScanProcFunc
	ThreadNum int
}

func (scmd *ScanCmd) clientScan(client *redis.Client) error {
	var scanPos uint64 = 0
	scanLock := sync.Mutex{}
	exitFlag := false

	// SCAN匹配关键字
	matchKey := scmd.Match
	if matchKey == "" {
		matchKey = "*"
	}
	// 每次扫描条数
	var scanCnt int64 = 5000
	if scmd.Count > 0 {
		scanCnt = scmd.Count
	}

	wg := &sync.WaitGroup{}
	for i := 0; i < scmd.ThreadNum; i++ {
		wg.Add(1)
		go func(index int) error {
			defer wg.Done()
			for {
				scanLock.Lock()
				if exitFlag {
					scanLock.Unlock()
					return nil
				}
				cmd := client.Scan(scanPos, matchKey, scanCnt)
				if cmd.Err() != nil {
					scanLock.Unlock()
					return cmd.Err()
				}
				keys, pos := cmd.Val()
				if pos > 0 {
					scanPos = pos
				} else {
					exitFlag = true
				}
				scanLock.Unlock()
				if scmd.ScanFunc != nil {
					for _, k := range keys {
						scmd.ScanFunc(k)
					}
				}
			}
			return nil
		}(i)
	}
	wg.Wait()
	return nil
}

func Scan(match string, count int64, scanFunc ScanProcFunc, threadNum int) error {
	cmd := &ScanCmd{
		Match:     match,
		Count:     count,
		ScanFunc:  scanFunc,
		ThreadNum: threadNum,
	}
	return DefaultCluster.client.ForEachMaster(cmd.clientScan)
}

func (c *RedisCluster) Scan(match string, count int64, scanFunc ScanProcFunc, threadNum int) error {
	cmd := &ScanCmd{
		Match:     match,
		Count:     count,
		ScanFunc:  scanFunc,
		ThreadNum: threadNum,
	}
	return c.client.ForEachMaster(cmd.clientScan)
}

func Hset(key, field string, value interface{}) (bool, error) {
	return boolCmdRes(DefaultCluster.client.HSet(key, field, value))
}

func (c *RedisCluster) Hset(key, field string, value interface{}) (bool, error) {
	return boolCmdRes(c.client.HSet(key, field, value))
}

func HGet(key, field string) (string, error) {
	return stringCmdRes(DefaultCluster.client.HGet(key, field))
}

func (c *RedisCluster) HGet(key, field string) (string, error) {
	return stringCmdRes(c.client.HGet(key, field))
}
