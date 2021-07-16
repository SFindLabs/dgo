package crypto

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"golang.org/x/crypto/pkcs12"
)

func FormatCertificate(raw string) []byte {
	return FormatKey(raw, "-----BEGIN CERTIFICATE-----", "-----END CERTIFICATE-----", 76)
}

func ParseCertificate(b []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(b)
	if block == nil {
		return nil, errors.New("certificate failed to load")
	}
	csr, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	return csr, nil
}

func Pkcs12ToPem(p12 []byte, password string) (cert tls.Certificate, err error) {
	blocks, err := pkcs12.ToPEM(p12, password)

	if err != nil {
		return cert, err
	}

	var pemData []byte
	for _, b := range blocks {
		pemData = append(pemData, pem.EncodeToMemory(b)...)
	}

	cert, err = tls.X509KeyPair(pemData, pemData)
	return cert, err
}
