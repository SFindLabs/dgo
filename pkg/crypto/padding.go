package crypto

import (
	"bytes"
)

func ZeroPadding(data []byte, blockSize int) []byte {
	diff := blockSize - len(data)%blockSize
	paddingText := bytes.Repeat([]byte{0}, diff)
	return append(data, paddingText...)
}

func PKCS7Padding(data []byte, blockSize int) []byte {
	diff := blockSize - len(data)%blockSize
	paddingText := bytes.Repeat([]byte{byte(diff)}, diff)
	return append(data, paddingText...)
}

func PKCS7UnPadding(data []byte) []byte {
	length := len(data)
	unPadding := int(data[length-1])
	return data[:(length - unPadding)]
}
