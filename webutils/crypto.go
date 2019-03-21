package webutils

import (
	"fmt"

	"github.com/wlibo666/common-lib/crypto"
)

func Md5Pwd(t int64, pwd string) string {
	return crypto.Md5([]byte(crypto.Md5([]byte(pwd)) + fmt.Sprintf("%d", t)))
}

func CheckPwd(t int64, md5Pwd, encryPwd string) bool {
	return crypto.Md5([]byte(md5Pwd+fmt.Sprintf("%d", t))) == encryPwd
}
