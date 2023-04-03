package encryption

import (
	"crypto/rsa"
	"encoding/base64"
	"hash"
	"io"
)

// EncryptionOpts is the options for encryption
type EncryptionOpts struct {
	Random    io.Reader
	Hash      hash.Hash
	PublicKey *rsa.PublicKey
	Label     []byte
}

// Encrypt will encrypt the data using rsa.EncryptOAEP
func Encrypt(data []byte, opts *EncryptionOpts) ([]byte, error) {
	enc, err := rsa.EncryptOAEP(opts.Hash, opts.Random, opts.PublicKey, data, opts.Label)
	if err != nil {
		return nil, err
	}

	return enc, nil
}

// EncryptToBase64 wrapper for Encrypt then encode the output to base64
func EncryptToBase64(data []byte, opts *EncryptionOpts) (string, error) {
	enc, err := Encrypt(data, opts)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(enc), nil
}

// DecryptionOpts is the options for decryption
type DecryptionOpts struct {
	PrivateKey *rsa.PrivateKey
	Random     io.Reader
	Hash       hash.Hash
	Label      []byte
}

// Decrypt will decrypt the data using rsa.DecryptOAEP
func Decrypt(data []byte, opts *DecryptionOpts) ([]byte, error) {
	decrypted, err := rsa.DecryptOAEP(opts.Hash, opts.Random, opts.PrivateKey, data, opts.Label)
	if err != nil {
		return nil, err
	}

	return decrypted, nil
}

// DecryptFromBase64 wrapper for Decrypt then decode the input from base64
func DecryptFromBase64(data string, opts *DecryptionOpts) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	return Decrypt(decoded, opts)
}
