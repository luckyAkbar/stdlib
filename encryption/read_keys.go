package encryption

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
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

// ReadPublicKey will read public key
func ReadPublicKey(key []byte) (*rsa.PublicKey, error) {
	public, _ := pem.Decode(key)
	if public == nil {
		return nil, errors.New("key error: unable to decode public key")
	}

	if public.Type != "PUBLIC KEY" {
		return nil, errors.New("key error: unknown type of public key")
	}

	publicKey, err := x509.ParsePKCS1PublicKey(public.Bytes)
	if err != nil {
		return nil, err
	}

	return publicKey, nil
}

// ReadPublicKeyFromFile wrapper for ReadPublicKey with option to read file based on path location
func ReadPublicKeyFromFile(path string) (*rsa.PublicKey, error) {
	r, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}

	return ReadPublicKey(r)
}
