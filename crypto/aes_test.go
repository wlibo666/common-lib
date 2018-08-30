package crypto

import (
	"encoding/base32"
	"encoding/base64"
	"testing"
)

func TestAesEncrypt(t *testing.T) {
	data, err := AesEncrypt([]byte("mypwd"), []byte("0123456789123456"))
	if err != nil {
		t.Fatalf("AesEncrypt failed,err:%s", err.Error())
	}
	t.Logf("base32 data:%s", base32.StdEncoding.EncodeToString(data))
	t.Logf("base64 data:%s", base64.StdEncoding.EncodeToString(data))

	text, err := AesDecrypt(data, []byte("0123456789123456"))
	if err != nil {
		t.Fatalf("AesDecrypt faileed,err:%s", err.Error())
	}
	t.Logf("text:%s", text)
}
