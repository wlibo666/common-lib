package log

import (
	"errors"

	"github.com/Sirupsen/logrus"
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
)

const (
	LOG_FORMAT_TEXT = iota
	LOG_FORMAT_JSON
)

var (
	DefFileLogger *logrus.Logger
)

func NewDefaultFileLogger(fileName string, fileCnt int, logFormat ...int) error {
	var err error
	DefFileLogger, err = NewFileLogger(fileName, fileCnt, logFormat...)
	if err != nil {
		return err
	}
	return nil
}

// 新建日志结构
// fileName: 日志文件名称
// fileCnt: 日志文件切分最大个数(每天一个)
// logFormat: 日志格式，文本或者json(LOG_FORMAT_TEXT|LOG_FORMAT_JSON)
func NewFileLogger(fileName string, fileCnt int, logFormat ...int) (*logrus.Logger, error) {
	if DefFileLogger != nil {
		return DefFileLogger, nil
	}
	if fileName == "" {
		return nil, errors.New("lost fileName")
	}
	if fileCnt <= 0 {
		fileCnt = 3
	}

	writer, err := rotatelogs.New(
		fileName+".%Y%m%d",
		rotatelogs.WithLinkName(fileName),
		rotatelogs.WithRotationCount(fileCnt),
	)
	if err != nil {
		return nil, err
	}

	pathMap := lfshook.WriterMap{}
	for _, level := range logrus.AllLevels {
		pathMap[level] = writer
	}
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
				ForceColors:      true,
				DisableColors:    true,
				DisableTimestamp: false,
			},
		)
	}

	logrus.AddHook(fileHook)

	logger := logrus.New()
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
