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

//---------------------------AES ECB---------------------------------

func AesEncryptECB(origData []byte, key []byte) (encrypted []byte, err error) {
	defer func() {
		if err := recover(); err != nil {
			err = errors.New(fmt.Sprintf("AesEncryptECB fail:%v", err))
		}
	}()
	block, err := aes.NewCipher(generateKey(key))
	if err != nil {
		return
	}
	length := (len(origData) + aes.BlockSize) / aes.BlockSize
	plain := make([]byte, length*aes.BlockSize)
	copy(plain, origData)
	pad := byte(len(plain) - len(origData))
	for i := len(origData); i < len(plain); i++ {
		plain[i] = pad
	}
	encrypted = make([]byte, len(plain))
	// 分组分块加密
	for bs, be := 0, block.BlockSize(); bs <= len(origData); bs, be = bs+block.BlockSize(), be+block.BlockSize() {
		block.Encrypt(encrypted[bs:be], plain[bs:be])
	}

	return encrypted, nil
}

func AesDecryptECB(encrypted []byte, key []byte) (decrypted []byte, err error) {
	defer func() {
		if err := recover(); err != nil {
			err = errors.New(fmt.Sprintf("AesDecryptECB fail:%v", err))
		}
	}()
	block, err := aes.NewCipher(generateKey(key))
	if err != nil {
		return
	}
	decrypted = make([]byte, len(encrypted))
	for bs, be := 0, block.BlockSize(); bs < len(encrypted); bs, be = bs+block.BlockSize(), be+block.BlockSize() {
		block.Decrypt(decrypted[bs:be], encrypted[bs:be])
	}

	trim := 0
	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}

	return decrypted[:trim], nil
}

func generateKey(key []byte) (genKey []byte) {
	genKey = make([]byte, 16)
	copy(genKey, key)
	for i := 16; i < len(key); {
		for j := 0; j < 16 && i < len(key); j, i = j+1, i+1 {
			genKey[j] ^= key[i]
		}
	}
	return genKey
}
