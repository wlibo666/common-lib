package rediscluster

import (
	"fmt"
	"testing"
	"time"
)

func TestSetNx(t *testing.T) {
	ClusterInit("10.110.92.171:6379,10.110.92.171:6380,10.110.92.171:6381,10.110.92.172:6379,10.110.92.172:6380,10.110.92.172:6381", "", 100)
	Del("mynx")
	succ, err := SetNx("mynx", "mynxval", 0)
	if err != nil {
		t.Fatalf("setnx failed,err:%v", err.Error())
	}
	fmt.Printf("setnx res:%v\n", succ)
	val, err := GetStr("mynx")
	if err != nil {
		t.Fatalf("get failed,err:%s", err.Error())
	}
	t.Logf("val:%s", val)
	succ, err = SetNx("mynx", "mynxval", 0)
	if err != nil {
		t.Fatalf("setnx 2 failed,err:%v", err.Error())
	}
	fmt.Printf("setnx 2 res:%v\n", succ)
}

func TestSet(t *testing.T) {
	ClusterInit("10.110.92.171:6379,10.110.92.171:6380,10.110.92.171:6381,10.110.92.172:6379,10.110.92.172:6380,10.110.92.172:6381", "", 100)
	Del("myset")
	succ, err := Set("myset", "mynxval", 0)
	if err != nil {
		t.Fatalf("set failed,err:%v", err.Error())
	}
	fmt.Printf("set res:%v\n", succ)
	val, err := GetStr("myset")
	if err != nil {
		t.Fatalf("get failed,err:%s", err.Error())
	}
	t.Logf("val:%s", val)
	succ, err = Set("myset", "mynxval", 0)
	if err != nil {
		t.Fatalf("setnx 2 failed,err:%v", err.Error())
	}
	fmt.Printf("set 2 res:%v\n", succ)
}

func TestDel(t *testing.T) {
	ClusterInit("10.110.92.171:6379,10.110.92.171:6380,10.110.92.171:6381,10.110.92.172:6379,10.110.92.172:6380,10.110.92.172:6381", "", 20)

	Del("program_calc_play_duration")
}

func TestHset(t *testing.T) {
	ClusterInit("10.110.92.171:6379,10.110.92.171:6380,10.110.92.171:6381,10.110.92.172:6379,10.110.92.172:6380,10.110.92.172:6381", "", 20)

	Hset("mytestkey", "testfield", true)
	v, err := HGet("mytestkey", "testfield")
	if err != nil {
		t.Fatalf("HGet failed,err:%s", err.Error())
	}
	t.Logf("value:%s", v)
}

func processKey(key string) error {
	fmt.Printf("process key:%s\n", key)
	return nil
}

func TestScan(t *testing.T) {
	ClusterInit("10.58.80.254:6379,10.58.80.254:6380,10.58.80.254:6381,10.58.80.254:6382,10.58.80.254:6383,10.58.80.254:6384", "", 20000)
	/*for i := 0; i < 5; i++ {
		k := fmt.Sprintf("k%d", i)
		v, _ := DefaultCluster.Del(k)
		if v == 1 {
			t.Logf("del k:%s ok", k)
		}
	}
	for i := 0; i < 5; i++ {
		k := fmt.Sprintf("k%d", i)
		v := fmt.Sprintf("v%d", i)
		res, _ := DefaultCluster.SetNx(k, v, 0)
		if res {
			t.Logf("set k:%s v:%s ok", k, v)
		}
	}
	for i := 0; i < 5; i++ {
		k := fmt.Sprintf("k%d", i)
		v, err := DefaultCluster.GetStr(k)
		if err != nil {
			t.Fatalf("get %s failed,err:%s", k, err.Error())
		}
		t.Logf("get k:%s,v:%s", k, v)
	}*/

	Scan("k*", 100, processKey, 2)
}

func TestExpire(t *testing.T) {
	ClusterInit("10.110.92.171:6379,10.110.92.171:6380,10.110.92.172:6379", "", 20000)

	SetNx("exkey", "exvalue", 0)
	Expire("exkey", 10)
	time.Sleep(3 * time.Second)
	ttl, err := TTL("exkey")
	if err != nil {
		t.Fatalf("ttl failed,err:%s", err.Error())
	}
	t.Logf("ttl is:%d", ttl)
	Del("exkey")
}

func TestGetInt64(t *testing.T) {
	ClusterInit("10.110.92.171:6379,10.110.92.171:6380,10.110.92.172:6379", "", 20000)

	SetNx("exkeyint", 10, 0)
	v, err := GetInt64("exkeyint")
	if err != nil {
		t.Fatalf("GetInt64 failed,err:%s", err.Error())
	}
	t.Logf("v is:%d", v)

	v, err = GetInt64("exkeyint222")
	if err != nil {
		t.Fatalf("GetInt64 failed,err:%s", err.Error())
	}
}

func TestIncrBy(t *testing.T) {
	ClusterInit("10.110.92.171:6379,10.110.92.171:6380,10.110.92.172:6379", "", 20000)

	SetNx("exkeyint", 10, 0)
	IncrBy("exkeyint", 5)
	v, err := GetInt64("exkeyint")
	if err != nil {
		t.Fatalf("GetInt64 failed,err:%s", err.Error())
	}
	t.Logf("v is:%d", v)

}

func TestZAdd(t *testing.T) {
	key := "testmysortset"
	ClusterInit("10.110.92.171:6379,10.110.92.171:6380,10.110.92.172:6379", "", 100)
	Del(key)

	ZAdd(key, float64(1514265354), 1514265354)
	ZAdd(key, float64(1514265355), 1514265355)
	ZAdd(key, float64(1514265350), 1514265350)
	ZAdd(key, float64(1514265359), 1514265359)
	ZAdd(key, float64(1514265344), 1514265344)

	vs, err := ZRange(key, 0, -1)
	if err != nil {
		t.Fatalf("zrange err:%s", err.Error())
	}
	t.Logf("vs:%v", vs)

	ZRemRangeByRank(key, 1, -2)
	vs, err = ZRange(key, 0, -1)
	if err != nil {
		t.Fatalf("zrange err:%s", err.Error())
	}
	t.Logf("vs:%v", vs)
	Del(key)
}
