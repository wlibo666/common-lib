package log

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func TestNewFileLogger(t *testing.T) {
	err := NewDefaultFileLogger("t.log", 2)
	if err != nil {
		t.Fatalf("new logger failed,err:%s", err.Error())
	}
	SetDefaultLoggerLevel(logrus.DebugLevel)

	for i := 0; i < 60; i++ {
		t.Logf("i is:%d\n", i)
		DefFileLogger.Debugf("i is:%d", i)
		DefFileLogger.WithFields(
			logrus.Fields{
				"for":    "withfield",
				"author": "wangchunyan",
				"cnt":    2,
			}).Debugf("i is:%d", i)
		time.Sleep(time.Second)
	}
}
