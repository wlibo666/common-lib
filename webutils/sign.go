package webutils

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/wlibo666/common-lib/crypto"
)

const (
	ARG_SIGN = "_sign"
)

var (
	PARAM_SIGN = ARG_SIGN
)

func SetParamSign(sign string) {
	PARAM_SIGN = sign
}

func GenRequestArgs(args map[string]string) string {
	var argString []string
	for _, k := range args {
		argString = append(argString, fmt.Sprintf("%s=%s", k, args[k]))
	}
	return strings.Join(argString, "&")
}

func GenSortRequestArgs(args map[string]string) string {
	var slice sort.StringSlice
	var argString []string
	if args == nil {
		return ""
	}
	for k, _ := range args {
		slice = append(slice, k)
	}
	sort.Sort(slice)
	for _, k := range slice {
		argString = append(argString, fmt.Sprintf("%s=%s", k, args[k]))
	}
	return strings.Join(argString, "&")
}

func GenSignature(method, requestURI string, args map[string]string, secret string) string {
	content := ""
	content += method
	content += "\r\n"
	content += requestURI
	content += "\r\n"
	if args != nil {
		content += GenSortRequestArgs(args)
		content += "\r\n"
	}
	content += secret
	return fmt.Sprintf("%x", md5.Sum([]byte(content)))
}

func Encode3DESBase64(data, key []byte) (string, error) {
	d, err := crypto.TripleDesEncrypt(data, key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(d), nil
}

func DecodeBase643DES(data, key []byte) ([]byte, error) {
	des, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return []byte{}, err
	}
	d, err := crypto.TripleDesDecrypt(des, key)
	if err != nil {
		return []byte{}, err
	}
	return d, nil
}

func genCharDigit() byte {
	rand.Seed(time.Now().UnixNano())
	c := rand.Int() % 123
	if c >= 97 || (c >= 65 && c <= 90) || (c >= 48 && c <= 57) {
		return byte(c)
	}
	return genCharDigit()
}

func GenRandomStringN(n int) string {
	var str []byte
	_n := n
	for {
		str = append(str, genCharDigit())
		_n--
		if _n <= 0 {
			break
		}
	}
	return string(str)
}
