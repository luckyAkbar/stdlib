package encryption

import (
	"crypto"
	"crypto/rsa"
	"encoding/base64"
	"io"
)

// SignOpts is a struct that contains the options for signing a message
type SignOpts struct {
	Random  io.Reader
	PrivKey *rsa.PrivateKey
	Alg     crypto.Hash
	PSSOpts *rsa.PSSOptions
}

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

func SignToBase64(message []byte, opts *SignOpts) (string, error) {
	signature, err := Sign(message, opts)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

type VerifyOpts struct {
	PublicKey *rsa.PublicKey
	Alg       crypto.Hash
	PSSOpts   *rsa.PSSOptions
}

func Verify(rawMessage, signature []byte, opts *VerifyOpts) error {
	msgHash := opts.Alg.New()
	_, err := msgHash.Write(rawMessage)
	if err != nil {
		return err
	}

	msgHashSum := msgHash.Sum(nil)

	return rsa.VerifyPSS(opts.PublicKey, opts.Alg, msgHashSum, signature, opts.PSSOpts)
}
