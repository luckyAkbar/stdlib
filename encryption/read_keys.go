package encryption

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"path/filepath"
)

// KeyComponent is a struct that contains private key, public key, and bytes
type KeyComponent struct {
	// PrivateKey will only be populated when reading private key
	PrivateKey *rsa.PrivateKey

	// PublicKey will always be populated
	PublicKey *rsa.PublicKey

	// Bytes contains actual bytes of key file
	Bytes []byte
}

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
func ReadKeyFromFile(path string) (*KeyComponent, error) {
	r, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}

	key, err := ReadKey(r)
	if err != nil {
		return nil, err
	}

	return &KeyComponent{
		PrivateKey: key,
		PublicKey:  &key.PublicKey,
		Bytes:      r,
	}, nil
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
// this function will return only public key and set private key to nil
func ReadPublicKeyFromFile(path string) (*KeyComponent, error) {
	r, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}

	key, err := ReadPublicKey(r)
	if err != nil {
		return nil, err
	}

	return &KeyComponent{
		PublicKey: key,
		Bytes:     r,
	}, nil
}
