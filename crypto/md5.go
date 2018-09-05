package crypto

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5(data []byte) string {
	hash := md5.Sum(data)
	return hex.EncodeToString([]byte(hash[:md5.Size]))
}
