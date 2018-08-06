package crypto

import (
	"encoding/base64"
	"testing"
)

var (
	key        = []byte("012345678901234567890123")
	clearText  = []byte("hello,world.")
	cipherText = "tdjwn3NgyFUKM7IlbCF09Q=="
)

func TestFillNBytes(t *testing.T) {
	newKey1 := FillNBytes([]byte("0123"), 24)
	t.Logf("newkey1:%s", newKey1)
	newKey2 := FillNBytes([]byte(""), 24)
	t.Logf("newkey1:%sEND", newKey2)
	newKey3 := FillNBytes([]byte("012345678901234567890123"), 24)
	t.Logf("newkey1:%s", newKey3)
	newKey4 := FillNBytes([]byte("012345678901234567890123456"), 24)
	t.Logf("newkey1:%s", newKey4)
}

func TestTripleDESEncrypt(t *testing.T) {
	dst, err := TripleDesEncrypt(clearText, key)
	if err != nil {
		t.Fatalf("TripleDesEncrypt failed,err:%s", err.Error())
	}
	t.Logf("dst:%x, %s", dst, base64.StdEncoding.EncodeToString(dst))

	if base64.StdEncoding.EncodeToString(dst) != cipherText {
		t.Fatalf("encrypto failed")
	}
}

func TestTripleDESDecrypt(t *testing.T) {
	cipher, _ := base64.StdEncoding.DecodeString(cipherText)
	t.Logf("%x, %s", cipher, cipherText)
	dst, err := TripleDesDecrypt(cipher, key)
	if err != nil {
		t.Fatalf("TripleDesDecrypt failed,err:%s", err.Error())
	}
	t.Logf("dst:%s", dst)
	if string(dst) != string(clearText) {
		t.Fatalf("TripleDesDecrypt failed")
	}
}