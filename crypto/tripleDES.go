package crypto

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
)

const (
	TRIPLE_DES_KEY_LEN = 24
)

func FillNBytes(key []byte, n int) []byte {
	var tmpKey []byte
	var tmpLen = len(key)
	if tmpLen >= n {
		return key[:n]
	}
	if tmpLen == 0 {
		key = append(key, ' ')
		tmpLen = 1
	}
	for i := 0; i < n; i++ {
		tmpKey = append(tmpKey, key[i%tmpLen])
	}
	return tmpKey
}

// 3DES加密
func TripleDesEncrypt(origData, key []byte) ([]byte, error) {
	key = FillNBytes(key, TRIPLE_DES_KEY_LEN)
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	origData = PKCS5Padding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key[:8])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// 3DES解密
func TripleDesDecrypt(crypted, key []byte) ([]byte, error) {
	key = FillNBytes(key, TRIPLE_DES_KEY_LEN)
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, key[:8])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
