package crypto

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"strings"
)

const (
	kPublicKeyPrefix = "-----BEGIN PUBLIC KEY-----"
	kPublicKeySuffix = "-----END PUBLIC KEY-----"

	kPKCS1Prefix = "-----BEGIN RSA PRIVATE KEY-----"
	KPKCS1Suffix = "-----END RSA PRIVATE KEY-----"

	kPKCS8Prefix = "-----BEGIN PRIVATE KEY-----"
	KPKCS8Suffix = "-----END PRIVATE KEY-----"

	kPublicKeyType     = "PUBLIC KEY"
	kPrivateKeyType    = "PRIVATE KEY"
	kRSAPrivateKeyType = "RSA PRIVATE KEY"
)

var (
	errPrivateKeyFailedToLoad = errors.New("private key failed to load")
	errPublicKeyFailedToLoad  = errors.New("public key failed to load")
)

func FormatPublicKey(raw string) []byte {
	return FormatKey(raw, kPublicKeyPrefix, kPublicKeySuffix, 64)
}

func FormatPKCS1PrivateKey(raw string) []byte {
	raw = strings.Replace(raw, kPKCS8Prefix, "", 1)
	raw = strings.Replace(raw, KPKCS8Suffix, "", 1)
	return FormatKey(raw, kPKCS1Prefix, KPKCS1Suffix, 64)
}

func FormatPKCS8PrivateKey(raw string) []byte {
	raw = strings.Replace(raw, kPKCS1Prefix, "", 1)
	raw = strings.Replace(raw, KPKCS1Suffix, "", 1)
	return FormatKey(raw, kPKCS8Prefix, KPKCS8Suffix, 64)
}

func ParsePKCS1PrivateKey(data []byte) (key *rsa.PrivateKey, err error) {
	var block *pem.Block
	block, _ = pem.Decode(data)
	if block == nil {
		return nil, errPrivateKeyFailedToLoad
	}

	key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return key, err
}

func ParsePKCS8PrivateKey(data []byte) (key *rsa.PrivateKey, err error) {
	var block *pem.Block
	block, _ = pem.Decode(data)
	if block == nil {
		return nil, errPrivateKeyFailedToLoad
	}

	rawKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	key, ok := rawKey.(*rsa.PrivateKey)
	if ok == false {
		return nil, errPrivateKeyFailedToLoad
	}

	return key, err
}

func ParsePublicKey(data []byte) (key *rsa.PublicKey, err error) {
	var block *pem.Block
	block, _ = pem.Decode(data)
	if block == nil {
		return nil, errPublicKeyFailedToLoad
	}

	var pubInterface interface{}
	pubInterface, err = x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	key, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return nil, errPublicKeyFailedToLoad
	}

	return key, err
}

func packageData(data []byte, packageSize int) (r [][]byte) {
	src := make([]byte, len(data))
	copy(src, data)

	r = make([][]byte, 0)
	if len(src) <= packageSize {
		return append(r, src)
	}
	for len(src) > 0 {
		p := src[:packageSize]
		r = append(r, p)
		src = src[packageSize:]
		if len(src) <= packageSize {
			r = append(r, src)
			break
		}
	}
	return r
}

// RSAEncrypt ???????????? key ????????? data ?????? RSA ??????
func RSAEncrypt(data, key []byte) ([]byte, error) {
	pubKey, err := ParsePublicKey(key)
	if err != nil {
		return nil, err
	}

	return RSAEncryptWithKey(data, pubKey)
}

// RSAEncryptWithKey ???????????? key ????????? data ?????? RSA ??????
func RSAEncryptWithKey(data []byte, key *rsa.PublicKey) ([]byte, error) {
	pData := packageData(data, key.N.BitLen()/8-11)
	cipher := make([]byte, 0)

	for _, d := range pData {
		c, e := rsa.EncryptPKCS1v15(rand.Reader, key, d)
		if e != nil {
			return nil, e
		}
		cipher = append(cipher, c...)
	}

	return cipher, nil
}

// RSADecryptWithPKCS1 ???????????? key ????????? data ?????? RSA ?????????key ???????????? pkcs1
func RSADecryptWithPKCS1(data, key []byte) ([]byte, error) {
	priKey, err := ParsePKCS1PrivateKey(key)
	if err != nil {
		return nil, err
	}

	return RSADecryptWithKey(data, priKey)
}

// RSADecryptWithPKCS1 ???????????? key ????????? data ?????? RSA ?????????key ???????????? pkcs8
func RSADecryptWithPKCS8(data, key []byte) ([]byte, error) {
	priKey, err := ParsePKCS8PrivateKey(key)
	if err != nil {
		return nil, err
	}

	return RSADecryptWithKey(data, priKey)
}

// RSADecryptWithKey ???????????? key ????????? data ?????? RSA ??????
func RSADecryptWithKey(data []byte, key *rsa.PrivateKey) ([]byte, error) {
	pData := packageData(data, key.PublicKey.N.BitLen()/8)
	plain := make([]byte, 0, 0)

	for _, d := range pData {
		p, e := rsa.DecryptPKCS1v15(rand.Reader, key, d)
		if e != nil {
			return nil, e
		}
		plain = append(plain, p...)
	}
	return plain, nil
}

func RSASignWithPKCS1(data, key []byte, hash crypto.Hash) ([]byte, error) {
	priKey, err := ParsePKCS1PrivateKey(key)
	if err != nil {
		return nil, err
	}
	return RSASignWithKey(data, priKey, hash)
}

func RSASignWithPKCS8(data, key []byte, hash crypto.Hash) ([]byte, error) {
	priKey, err := ParsePKCS8PrivateKey(key)
	if err != nil {
		return nil, err
	}
	return RSASignWithKey(data, priKey, hash)
}

func RSASignWithKey(data []byte, key *rsa.PrivateKey, hash crypto.Hash) ([]byte, error) {
	h := hash.New()
	h.Write(data)
	hashed := h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, key, hash, hashed)
}

func RSAVerify(data, sig, key []byte, hash crypto.Hash) error {
	pubKey, err := ParsePublicKey(key)
	if err != nil {
		return err
	}
	return RSAVerifyWithKey(data, sig, pubKey, hash)
}

func RSAVerifyWithKey(data, sig []byte, key *rsa.PublicKey, hash crypto.Hash) error {
	h := hash.New()
	h.Write(data)
	hashed := h.Sum(nil)
	return rsa.VerifyPKCS1v15(key, hash, hashed, sig)
}

func getPublicKeyBytes(publicKey *rsa.PublicKey) ([]byte, error) {
	pubDer, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, err
	}

	pubBlock := &pem.Block{Type: kPublicKeyType, Bytes: pubDer}

	var pubBuf bytes.Buffer
	if err = pem.Encode(&pubBuf, pubBlock); err != nil {
		return nil, err
	}
	return pubBuf.Bytes(), nil
}

func GenRSAKeyWithPKCS1(bits int) (privateKey, publicKey []byte, err error) {
	priKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	priDer := x509.MarshalPKCS1PrivateKey(priKey)
	priBlock := &pem.Block{Type: kRSAPrivateKeyType, Bytes: priDer}

	var priBuf bytes.Buffer
	if err = pem.Encode(&priBuf, priBlock); err != nil {
		return nil, nil, err
	}

	publicKey, err = getPublicKeyBytes(&priKey.PublicKey)
	if err != nil {
		return nil, nil, err
	}
	privateKey = priBuf.Bytes()
	return privateKey, publicKey, err
}

func GenRSAKeyWithPKCS8(bits int) (privateKey, publicKey []byte, err error) {
	priKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	priDer, err := x509.MarshalPKCS8PrivateKey(priKey)
	if err != nil {
		return nil, nil, err
	}
	priBlock := &pem.Block{Type: kPrivateKeyType, Bytes: priDer}

	var priBuf bytes.Buffer
	if err = pem.Encode(&priBuf, priBlock); err != nil {
		return nil, nil, err
	}

	publicKey, err = getPublicKeyBytes(&priKey.PublicKey)
	if err != nil {
		return nil, nil, err
	}
	privateKey = priBuf.Bytes()

	return privateKey, publicKey, err
}
