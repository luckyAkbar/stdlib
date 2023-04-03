package encryption

import (
	"crypto"
	"crypto/rsa"
	"encoding/base64"
	"io"
)

// SignOpts is a struct that contains the options for signing a message.
// all struct fields are required unless otherwise noted.
type SignOpts struct {
	Random  io.Reader
	PrivKey *rsa.PrivateKey
	Alg     crypto.Hash
	PSSOpts *rsa.PSSOptions
}

// Sign will generate signature based on supplied message
func Sign(message []byte, opts *SignOpts) ([]byte, error) {
	msgHash := opts.Alg.New()
	_, err := msgHash.Write(message)
	if err != nil {
		return nil, err
	}

	msgHashSum := msgHash.Sum(nil)

	signature, err := rsa.SignPSS(opts.Random, opts.PrivKey, opts.Alg, msgHashSum, opts.PSSOpts)
	if err != nil {
		return nil, err
	}

	return signature, nil
}

// SignToBase64 wrapper for Sign with the output are base64 encoded string
func SignToBase64(message []byte, opts *SignOpts) (string, error) {
	signature, err := Sign(message, opts)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

// VerifyOpts is a struct that contains the options for verifying a message.
// all struct fields are required unless otherwise noted.
type VerifyOpts struct {
	PublicKey *rsa.PublicKey
	Alg       crypto.Hash
	PSSOpts   *rsa.PSSOptions
}

// Verify will verify the signature of a message.
// a valid signature will return nil, otherwise a non-nil error returned
func Verify(rawMessage, signature []byte, opts *VerifyOpts) error {
	msgHash := opts.Alg.New()
	_, err := msgHash.Write(rawMessage)
	if err != nil {
		return err
	}

	msgHashSum := msgHash.Sum(nil)

	return rsa.VerifyPSS(opts.PublicKey, opts.Alg, msgHashSum, signature, opts.PSSOpts)
}
