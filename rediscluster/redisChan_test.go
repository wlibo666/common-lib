package rediscluster

import (
	"testing"
	"time"
)

func TestSendRedisCmd(t *testing.T) {
	ClusterInit("10.110.92.171:6379,10.110.92.171:6380,10.110.92.171:6381,10.110.92.172:6379,10.110.92.172:6380,10.110.92.172:6381", "", 10)
	StartRedisChan()

	Del("mycmdstr")
	Del("mycmdstrnx")
	Del("mycmdhset")

	SendRedisCmd(CMD_SETSTR, "mycmdstr", "myvalue")
	SendRedisCmd(CMD_SETNXSTR, "mycmdstrnx", "mysetnxvalue")
	SendRedisCmd(CMD_HSETSTR, "mycmdhset", "myfiled", "myhsetvalue")
	time.Sleep(1 * time.Second)

	v, err := GetStr("mycmdstr")
	if err != nil {
		t.Fatalf("get failed,err:%s", err.Error())
	}
	t.Logf("v:%s", v)
	Del("mycmdstr")

	v, err = GetStr("mycmdstrnx")
	if err != nil {
		t.Fatalf("get failed,err:%s", err.Error())
	}
	t.Logf("v:%s", v)
	Del("mycmdstrnx")

	v, err = HGet("mycmdhset", "myfiled")
	if err != nil {
		t.Fatalf("get failed,err:%s", err.Error())
	}
	t.Logf("v:%s", v)
	Del("mycmdhset")

}
