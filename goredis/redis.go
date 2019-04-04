package goredis

import (
	"sync"

	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

var (
	redisPools = &sync.Map{}

	ERR_NOT_FOUND_CLIENT = errors.New("not found client in redis pools")
)

const (
	DEFAULT_NAME = "default"
)

func InitClient(addr, pwd string, db int, name ...string) error {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd,
		DB:       db,
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

func SetString(key, value string, expirySeconds int, name ...string) error {
	cli := GetClient(name...)
	if cli == nil {
		return ERR_NOT_FOUND_CLIENT
	}
	cmd := cli.Set(key, value, time.Duration(expirySeconds)*time.Second)
	_, err := cmd.Result()
	return err
}
