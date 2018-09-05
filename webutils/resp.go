package webutils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CommonResp struct {
	Errno int         `json:"errno"`
	Err   string      `json:"err"`
	Data  interface{} `json:"data"`
}

type NoneData struct {
}

type AddrData struct {
	Addr string `json:"addr"`
}

func ServeError(errno int, err string, ctx *gin.Context) {
	resp := CommonResp{
		Errno: errno,
		Err:   err,
		Data:  NoneData{},
	}
	ctx.JSON(http.StatusOK, resp)
	ctx.Abort()
}

func ServeResp(data interface{}, ctx *gin.Context) {
	resp := CommonResp{
		Errno: ERRNO_SUCCESS,
		Err:   "",
		Data:  data,
	}
	ctx.JSON(http.StatusOK, resp)
}

func IsValidResp(code int, err error) bool {
	if err != nil || code != http.StatusOK {
		return false
	}
	return true
}
