package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

//cbc 加密, iv等于16位key
func CbcEncrypt(encodeStr string, key []byte) (b string, rerr error) {
	defer func() {
		if err := recover(); err != nil {
			rerr = errors.New(fmt.Sprintf("AesEncrypt fail:%v", err))
		}
	}()

	encodeBytes := []byte(encodeStr)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	encodeBytes = zeroPadding(encodeBytes, blockSize)

	blockMode := cipher.NewCBCEncrypter(block, key)
	crypt := make([]byte, len(encodeBytes))
	blockMode.CryptBlocks(crypt, encodeBytes)

	return base64.StdEncoding.EncodeToString(crypt), nil
}

//cbc 解密, iv等于16位key
func CbcDecrypt(decodeStr string, key []byte) (b string, rerr error) {
	defer func() {
		if err := recover(); err != nil {
			rerr = errors.New(fmt.Sprintf("AesDecrypt fail:%v", err))
		}
	}()

	//先解密base64
	decodeBytes, err := base64.StdEncoding.DecodeString(decodeStr)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	blockMode := cipher.NewCBCDecrypter(block, key)
	origData := make([]byte, len(decodeBytes))

	blockMode.CryptBlocks(origData, decodeBytes)
	origData, err = zeroUnPadding(origData)
	if err != nil {
		return "", err
	}
	toStrData := strings.Trim(string(origData), "\x00")
	return toStrData, nil
}

func zeroPadding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{0}, padding)
	return append(cipherText, padText...)
}

func zeroUnPadding(origData []byte) ([]byte, error) {
	length := len(origData)
	unPadding := int(origData[length-1])

	if unPadding < 0 || unPadding > length {
		return nil, errors.New("may key error")
	}

	return origData[:(length - unPadding)], nil
}

func UrlBaseEncode(str string) string {
	str = strings.Replace(str, "+", "-", -1)
	str = strings.Replace(str, "/", "_", -1)
	str = strings.Replace(str, "=", "", -1)
	return str
}

func UrlBaseDecode(str string) string {
	var buffer bytes.Buffer
	str = strings.Replace(str, "-", "+", -1)
	str = strings.Replace(str, "_", "/", -1)
	buffer.WriteString(str)
	count := 4 - len(str)%4
	if count != 4 {
		for i := 0; i < count; i++ {
			buffer.WriteString("=")
		}
	}
	return buffer.String()
}
