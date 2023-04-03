package encryption

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
)

// ReadKey will read private key in form of []byte
func ReadKey(key []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(key)

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

// ReadKeyFromFile wrapper for ReadKey with option to read file based on path location
func ReadKeyFromFile(path string) (*rsa.PrivateKey, error) {
	r, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}

	return ReadKey(r)
}
