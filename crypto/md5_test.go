package crypto

import "testing"

func TestMd5(t *testing.T) {
	var key string = "hello,world."
	var _t string = "1535698092"
	t.Logf("%s md5:%s", key, Md5([]byte(key)))
	t.Logf("crypto:%s", Md5([]byte(string(Md5([]byte(key)))+_t)))
}
