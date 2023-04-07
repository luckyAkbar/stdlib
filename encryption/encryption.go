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

// EncryptWithSteps will encrypt the data using rsa.EncryptOAEP chunk by chunk. This is useful when the data is
// too large to be encrypted using Encrypt
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

// FileEncryptionOpts is the options for file encryption & decryption
type FileEncryptionOpts struct {
	SourcePath   string
	OutputPath   string
	AESKeyLength AESKeyLength
	Key          *KeyComponent

	// BufferSize size must be multiple of 16 bytes
	BufferSize int
}

// GetKeyLength will return the key length based on AESKeyLength
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

// GetChiperKey will return the key used to create chiper block
func (feo *FileEncryptionOpts) GetChiperKey() []byte {
	return feo.Key.Bytes[:feo.GetKeyLength()]
}

// EncryptFile will encrypt the file using AES
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

// DecryptFile will decrypt the file using AES and save the decrypted file to opts.OutputPath.
// this function will not delete the decrypted file, so it's up to the caller to delete the file after use.
func DecryptFile(opts *FileEncryptionOpts) error {
	infile, err := os.Open(opts.SourcePath)
	if err != nil {
		return nil, &custerr.ErrChain{
			Message: "unable to open source file for decryption",
			Cause:   err,
			Code:    http.StatusInternalServerError,
		}
	}
	defer infile.Close()

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
	fi, err := infile.Stat()
	if err != nil {
		return nil, &custerr.ErrChain{
			Message: "failed to get file info",
			Cause:   err,
			Code:    http.StatusInternalServerError,
		}
	}

	iv := make([]byte, block.BlockSize())
	msgLen := fi.Size() - int64(len(iv))
	_, err = infile.ReadAt(iv, msgLen)
	if err != nil {
		return nil, &custerr.ErrChain{
			Message: "failed to read file's chunks",
			Cause:   err,
			Code:    http.StatusInternalServerError,
		}
	}

	outfile, err := os.OpenFile(opts.OutputPath, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return nil, &custerr.ErrChain{
			Message: "failed open destination decrypted file",
			Cause:   err,
			Code:    http.StatusInternalServerError,
		}
	}
	defer helper.WrapCloser(outfile.Close)

	// The buffer size must be multiple of 16 bytes
	buf := make([]byte, opts.BufferSize)
	stream := cipher.NewCTR(block, iv)
	for {
		n, err := infile.Read(buf)
		if n > 0 {
			// The last bytes are the IV, don't belong the original message
			if n > int(msgLen) {
				n = int(msgLen)
			}
			msgLen -= int64(n)
			stream.XORKeyStream(buf, buf[:n])
			// Write into file
			outfile.Write(buf[:n])
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			break
		}
	}

	decrypted, err := os.ReadFile(opts.OutputPath)
	if err != nil {
		return nil, &custerr.ErrChain{
			Message: "failed to read decrypted file",
			Cause:   err,
			Code:    http.StatusInternalServerError,
		}
	}

	defer helper.DeleteFile(opts.OutputPath)

	return decrypted, nil
}
