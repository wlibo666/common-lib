package log

import (
	"errors"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

const (
	LOG_FORMAT_TEXT = iota
	LOG_FORMAT_JSON
)

var (
	DefFileLogger *logrus.Logger = &logrus.Logger{
		Out:       os.Stdout,
		Formatter: new(logrus.JSONFormatter),
		Level:     logrus.DebugLevel,
	}
)

func NewDefaultFileLogger(fileName string, fileCnt int, logFormat ...int) error {
	var err error
	DefFileLogger, err = NewFileLogger(fileName, fileCnt, logFormat...)
	if err != nil {
		return err
	}
	return nil
}

type DevNullWriter struct {
}

func (w *DevNullWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

// 新建日志结构
// fileName: 日志文件名称
// fileCnt: 日志文件切分最大个数(每天一个)
// logFormat: 日志格式，文本或者json(LOG_FORMAT_TEXT|LOG_FORMAT_JSON)
func NewFileLogger(fileName string, fileCnt int, logFormat ...int) (*logrus.Logger, error) {
	if fileName == "" {
		return nil, errors.New("lost fileName")
	}
	if fileCnt <= 0 {
		fileCnt = 3
	}

	// debug log
	debugWriter, err := rotatelogs.New(
		fileName+".debug.%Y%m%d",
		rotatelogs.WithLinkName(fileName+".debug"),
		rotatelogs.WithMaxAge(time.Duration(fileCnt)*24*time.Hour),
		rotatelogs.WithRotationCount(fileCnt),
	)
	if err != nil {
		return nil, err
	}
	// info log
	infoWriter, err := rotatelogs.New(
		fileName+".info.%Y%m%d",
		rotatelogs.WithLinkName(fileName+".info"),
		rotatelogs.WithMaxAge(time.Duration(fileCnt)*24*time.Hour),
		rotatelogs.WithRotationCount(fileCnt),
	)
	if err != nil {
		return nil, err
	}
	// warn/error
	warnWriter, err := rotatelogs.New(
		fileName+".warn.%Y%m%d",
		rotatelogs.WithLinkName(fileName+".warn"),
		rotatelogs.WithMaxAge(time.Duration(fileCnt)*24*time.Hour),
		rotatelogs.WithRotationCount(fileCnt),
	)
	if err != nil {
		return nil, err
	}

	pathMap := lfshook.WriterMap{}
	pathMap[logrus.DebugLevel] = debugWriter
	pathMap[logrus.InfoLevel] = infoWriter
	pathMap[logrus.WarnLevel] = warnWriter
	pathMap[logrus.ErrorLevel] = warnWriter
	pathMap[logrus.FatalLevel] = warnWriter
	pathMap[logrus.PanicLevel] = warnWriter

	var fileHook *lfshook.LfsHook
	if len(logFormat) == 1 && logFormat[0] == LOG_FORMAT_JSON {
		fileHook = lfshook.NewHook(
			pathMap,
			&logrus.JSONFormatter{},
		)
	} else {
		fileHook = lfshook.NewHook(
			pathMap,
			&logrus.TextFormatter{
				DisableColors:    true,
				DisableTimestamp: false,
			},
		)
	}

	logrus.AddHook(fileHook)
	logger := &logrus.Logger{
		Out:       &DevNullWriter{},
		Formatter: new(logrus.JSONFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.DebugLevel,
	}
	logger.Hooks.Add(fileHook)

	return logger, nil
}

func SetDefaultLoggerLevel(level logrus.Level) {
	SetLoggerLevel(DefFileLogger, level)
}

func SetLoggerLevel(logger *logrus.Logger, level logrus.Level) {
	if logger == nil {
		return
	}
	logger.Level = level
}

func MiddleAccessLog(ctx *gin.Context) {
	start := time.Now()
	ctx.Next()
	end := time.Now()
	cost := end.Sub(start)

	DefFileLogger.Debugf("resp:[%d] cost:[%13v] client:[%15s] method:[%s] uri:[%s] user-agent:[%s] resp-len:[%d]",
		ctx.Writer.Status(), cost, ctx.ClientIP(), ctx.Request.Method,
		ctx.Request.RequestURI, ctx.Request.UserAgent(), ctx.Writer.Size())
}
