package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"hash"
	"io"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
	custerr "github.com/sweet-go/stdlib/error"
	"github.com/sweet-go/stdlib/helper"
)

// Opts is the options for encryption
type Opts struct {
	Random    io.Reader
	Hash      hash.Hash
	PublicKey *rsa.PublicKey
	Label     []byte
}

// Encrypt will encrypt the data using rsa.EncryptOAEP
func Encrypt(data []byte, opts *Opts) ([]byte, error) {
	enc, err := rsa.EncryptOAEP(opts.Hash, opts.Random, opts.PublicKey, data, opts.Label)
	if err != nil {
		return nil, err
	}

	return enc, nil
}

func EncryptWithSteps(data []byte, opts *Opts) ([]byte, error) {
	msgLen := len(data)
	step := opts.PublicKey.Size() - 2*opts.Hash.Size() - 2
	var encryptedBytes []byte

	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}

		encryptedBlockBytes, err := Encrypt(data[start:finish], opts)
		if err != nil {
			return nil, err
		}

		encryptedBytes = append(encryptedBytes, encryptedBlockBytes...)
	}

	return encryptedBytes, nil
}

// EncryptToBase64 wrapper for Encrypt then encode the output to base64
func EncryptToBase64(data []byte, opts *Opts) (string, error) {
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

type FileEncryptionOpts struct {
	SourcePath   string
	OutputPath   string
	AESKeyLength AESKeyLength
	Key          *KeyComponent

	// BufferSize size must be multiple of 16 bytes
	BufferSize int
}

func (feo *FileEncryptionOpts) GetKeyLength() int {
	switch feo.AESKeyLength {
	case AES128:
		return 16
	case AES192:
		return 24
	case AES256:
		return 32
	default:
		return 16
	}
}

func (feo *FileEncryptionOpts) GetKey() []byte {
	return feo.Key.Bytes[:feo.GetKeyLength()]
}

func EncryptFile(opts *FileEncryptionOpts) (iv []byte, err error) {
	source, err := os.Open(opts.SourcePath)
	if err != nil {
		return nil, &custerr.ErrChain{
			Message: "unable to open source file",
			Cause:   err,
			Code:    http.StatusInternalServerError,
		}
	}

	defer helper.WrapCloser(source.Close)

	block, err := aes.NewCipher(opts.GetKey())
	if err != nil {
		return nil, &custerr.ErrChain{
			Message: "failed to create chiper block",
			Cause:   err,
			Code:    http.StatusInternalServerError,
		}
	}

	// Never use more than 2^32 random nonces with a given key
	// because of the risk of repeat.
	iv = make([]byte, block.BlockSize())
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, &custerr.ErrChain{
			Message: "failed to read file",
			Cause:   err,
			Code:    http.StatusInternalServerError,
		}
	}

	dest, err := os.Create(opts.OutputPath)
	if err != nil {
		return nil, &custerr.ErrChain{
			Message: "unable to create destination file",
			Cause:   err,
			Code:    http.StatusInternalServerError,
		}
	}

	defer helper.WrapCloser(dest.Close)

	buf := make([]byte, opts.BufferSize)
	stream := cipher.NewCTR(block, iv)
	for {
		n, err := source.Read(buf)
		if n > 0 {
			stream.XORKeyStream(buf, buf[:n])
			// Write into file
			dest.Write(buf[:n])
		}

		if err == io.EOF {
			break
		}

		// should we report error?
		if err != nil {
			logrus.Warn("error while reading file to encryption", err)
			break
		}
	}

	dest.Write(iv)

	return iv, nil
}
