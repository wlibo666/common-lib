package webutils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/wlibo666/common-lib/logrus"
	"github.com/wlibo666/common-lib/utils"
)

func GetPostFormString(ctx *gin.Context, key string) (string, int) {
	value, ok := ctx.GetPostForm(key)
	if !ok {
		log.DefFileLogger.WithFields(logrus.Fields{
			ERR_FIELD_POSITION: utils.GetFileAndLine(),
			"key":              key,
		}).Warn("not found form string key.")
		return "", ERRNO_LOST_FORM_PARAM
	}
	return value, ERRNO_SUCCESS
}

func GetPostFormInt(ctx *gin.Context, key string) (int, int) {
	value, ok := ctx.GetPostForm(key)
	if !ok {
		log.DefFileLogger.WithFields(logrus.Fields{
			ERR_FIELD_POSITION: utils.GetFileAndLine(),
			"key":              key,
		}).Warn("not found form int key.")
		return 0, ERRNO_LOST_FORM_PARAM
	}
	n, err := strconv.Atoi(value)
	if err != nil {
		log.DefFileLogger.WithFields(logrus.Fields{
			ERR_FIELD_POSITION: utils.GetFileAndLine(),
			"key":              key,
			"value":            value,
		}).Warn("invalid form int param.")
		return 0, ERRNO_INVALID_FORM_PARAM
	}
	return n, ERRNO_SUCCESS
}

func GetURLParamString(ctx *gin.Context, key string) (string, int) {
	value := ctx.Param(key)
	if value == "" {
		log.DefFileLogger.WithFields(logrus.Fields{
			ERR_FIELD_POSITION: utils.GetFileAndLine(),
			"key":              key,
		}).Warn("not found url string param.")
		return "", ERRNO_LOST_URL_PARAM
	}
	return value, ERRNO_SUCCESS
}

func GetURLParamInt(ctx *gin.Context, key string) (int, int) {
	value := ctx.Param(key)
	if value == "" {
		log.DefFileLogger.WithFields(logrus.Fields{
			ERR_FIELD_POSITION: utils.GetFileAndLine(),
			"key":              key,
		}).Warn("not found url int param.")
		return 0, ERRNO_LOST_URL_PARAM
	}
	n, err := strconv.Atoi(value)
	if err != nil {
		log.DefFileLogger.WithFields(logrus.Fields{
			ERR_FIELD_POSITION: utils.GetFileAndLine(),
			"key":              key,
			"value":            value,
		}).Warn("invalid url int param.")
		return 0, ERRNO_INVALID_URL_PARAM
	}
	return n, ERRNO_SUCCESS
}

func GetQueryString(ctx *gin.Context, key string) (string, int) {
	value := ctx.Query(key)
	if value == "" {
		log.DefFileLogger.WithFields(logrus.Fields{
			ERR_FIELD_POSITION: utils.GetFileAndLine(),
			"key":              key,
		}).Warn("not found query string key.")
		return "", ERRNO_LOST_REQ_PARAM
	}
	return value, ERRNO_SUCCESS
}

func GetQueryInt(ctx *gin.Context, key string) (int, int) {
	value := ctx.Query(key)
	if value == "" {
		log.DefFileLogger.WithFields(logrus.Fields{
			ERR_FIELD_POSITION: utils.GetFileAndLine(),
			"key":              key,
		}).Warn("not found query int key.")
		return 0, ERRNO_LOST_REQ_PARAM
	}
	n, err := strconv.Atoi(value)
	if err != nil {
		log.DefFileLogger.WithFields(logrus.Fields{
			ERR_FIELD_POSITION: utils.GetFileAndLine(),
			"key":              key,
			"value":            value,
		}).Warn("invalid query int param.")
		return 0, ERRNO_INVALID_REQ_PARAM
	}
	return n, ERRNO_SUCCESS
}

func GetQueryInt64(ctx *gin.Context, key string) (int64, int) {
	value := ctx.Query(key)
	if value == "" {
		log.DefFileLogger.WithFields(logrus.Fields{
			ERR_FIELD_POSITION: utils.GetFileAndLine(),
			"key":              key,
		}).Warn("not found query int64 key.")
		return 0, ERRNO_LOST_REQ_PARAM
	}
	n, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		log.DefFileLogger.WithFields(logrus.Fields{
			ERR_FIELD_POSITION: utils.GetFileAndLine(),
			"key":              key,
			"value":            value,
		}).Warn("invalid query int64 param.")
		return 0, ERRNO_INVALID_REQ_PARAM
	}
	return n, ERRNO_SUCCESS
}

func GenPostArgs(args map[string]string) string {
	var arg []string
	for k, v := range args {
		tmp := fmt.Sprintf("%s=%s", k, v)
		arg = append(arg, tmp)
	}
	return strings.Join(arg, "&")
}
