package log

import (
	"testing"

	"github.com/Sirupsen/logrus"
)

func TestNewFileLogger(t *testing.T) {
	err := NewDefaultFileLogger("t.log", 2)
	if err != nil {
		t.Fatalf("new logger failed,err:%s", err.Error())
	}
	SetDefaultLoggerLevel(logrus.DebugLevel)

	DefFileLogger.Debugf("this is for:%s", "debugf")

	DefFileLogger.WithFields(
		logrus.Fields{
			"for":    "withfield",
			"author": "wangchunyan",
			"cnt":    3,
		}).Debugf("this is for:%s", "fieldDebugf")
}
