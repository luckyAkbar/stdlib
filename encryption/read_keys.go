package encryption

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

func ReadKey(key []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(key)

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}
