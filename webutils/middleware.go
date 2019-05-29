package webutils

import (
	"math"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/wlibo666/common-lib/logrus"
	"github.com/wlibo666/common-lib/utils"
)

const (
	ARG_T = "_t"
)

var (
	PARAM_T = ARG_T

	checkTimeRange = 60
)

func SetParamT(t string) {
	PARAM_T = t
}

func SetTimeRange(t int) {
	checkTimeRange = t
}

func CheckTime(ctx *gin.Context) {
	t, code := GetQueryInt64(ctx, PARAM_T)
	if code != ERRNO_SUCCESS {
		ServeError(code, GetErrno(code), ctx)
		return
	}

	now := time.Now().Unix()
	if math.Abs(float64(now-t)) > float64(checkTimeRange)*1.0 {
		log.DefFileLogger.WithFields(logrus.Fields{
			ERR_FIELD_POSITION: utils.GetFileAndLine(),
			"now":              now,
			PARAM_T:            t,
			"checkTimeRange":   checkTimeRange,
		}).Warn("Invalid param _t,time range is too large than config")
		ServeError(ERRNO_INVALID_REQ_PARAM, GetErrno(ERRNO_INVALID_REQ_PARAM), ctx)
		return
	}
}

var (
	_secret string
)

func SetSignSecret(secret string) {
	_secret = secret
}

func CheckSign(ctx *gin.Context) {
	// 校验签名
	pass := CheckSignByRequest(_secret, ctx.Request)
	if !pass {
		log.DefFileLogger.WithFields(logrus.Fields{
			ERR_FIELD_POSITION: utils.GetFileAndLine(),
		}).Warn("Invalid signature")
		ServeError(ERRNO_INVALID_SIGN, GetErrno(ERRNO_INVALID_SIGN), ctx)
		return
	}
}

func AccessLog(ctx *gin.Context) {
	log.MiddleAccessLog(ctx)
}
