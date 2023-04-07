package encryption_test

import (
	"crypto"
	"crypto/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sweet-go/stdlib/encryption"
	"github.com/sweet-go/stdlib/helper"
)

func TestEncrypt(t *testing.T) {
	key, err := encryption.GenerateKey(&encryption.KeyGenerationOpts{
		Random: rand.Reader,
		Bits:   2048,
	})
	assert.NoError(t, err)

	t.Run("ok", func(t *testing.T) {
		data := []byte("hello world")
		opts := &encryption.Opts{
			Random:    rand.Reader,
			Hash:      crypto.SHA512.New(),
			PublicKey: &key.PublicKey,
		}

		enc, err := encryption.Encrypt(data, opts)
		assert.NoError(t, err)

		decrOpts := &encryption.DecryptionOpts{
			PrivateKey: key,
			Hash:       crypto.SHA512.New(),
		}

		decrypted, err := encryption.Decrypt(enc, decrOpts)
		assert.NoError(t, err)

		assert.Equal(t, decrypted, data)
	})

	t.Run("msg too long", func(t *testing.T) {
		data := make([]byte, 999999999)
		opts := &encryption.Opts{
			Random:    rand.Reader,
			Hash:      crypto.SHA512.New(),
			PublicKey: &key.PublicKey,
		}

		_, err := encryption.Encrypt(data, opts)
		assert.Error(t, err)
	})

	t.Run("ok - encrypt to base64", func(t *testing.T) {
		data := []byte("hello world")
		opts := &encryption.Opts{
			Random:    rand.Reader,
			Hash:      crypto.SHA512.New(),
			PublicKey: &key.PublicKey,
		}

		enc, err := encryption.EncryptToBase64(data, opts)
		assert.NoError(t, err)

		decrOpts := &encryption.DecryptionOpts{
			PrivateKey: key,
			Hash:       crypto.SHA512.New(),
		}

		decrypted, err := encryption.DecryptFromBase64(enc, decrOpts)
		assert.NoError(t, err)

		assert.Equal(t, decrypted, data)
	})

	t.Run("decrypt error", func(t *testing.T) {
		decrOpts := &encryption.DecryptionOpts{
			PrivateKey: key,
			Hash:       crypto.SHA512.New(),
		}

		_, err = encryption.Decrypt([]byte("oke oke adkaoksdposdpaodpapo"), decrOpts)
		assert.Error(t, err)
	})

	t.Run("decrypt error - invalid base64", func(t *testing.T) {
		decrOpts := &encryption.DecryptionOpts{
			PrivateKey: key,
			Hash:       crypto.SHA512.New(),
		}

		_, err = encryption.DecryptFromBase64("awkwakwakwkawka", decrOpts)
		assert.Error(t, err)
	})
}

func TestEncryptWithSteps(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		// approx 1mb
		data := make([]byte, 1000000)
		key, err := encryption.GenerateKey(nil)
		assert.NoError(t, err)
		opts := &encryption.Opts{
			Random:    rand.Reader,
			Hash:      crypto.SHA512.New(),
			PublicKey: &key.PublicKey,
		}

		_, err = encryption.EncryptWithSteps(data, opts)
		assert.NoError(t, err)
	})
}

func TestFileEncryptionOpts(t *testing.T) {
	key, err := encryption.ReadKeyFromFile("./testdata/test_private.pem")
	assert.NoError(t, err)

	t.Run("ok get key length", func(t *testing.T) {
		opts := &encryption.FileEncryptionOpts{
			SourcePath:   "./testdata/txt_to_encrypt.txt",
			OutputPath:   "./testdata/encrypted.txt",
			AESKeyLength: encryption.AES128,
			Key:          key,
			BufferSize:   1024,
		}

		assert.Equal(t, 16, opts.GetKeyLength())
	})

	t.Run("ok get key length", func(t *testing.T) {
		opts := &encryption.FileEncryptionOpts{
			SourcePath:   "./testdata/txt_to_encrypt.txt",
			OutputPath:   "./testdata/encrypted.txt",
			AESKeyLength: encryption.AES192,
			Key:          key,
			BufferSize:   1024,
		}

		assert.Equal(t, 24, opts.GetKeyLength())
	})

	t.Run("ok get key length", func(t *testing.T) {
		opts := &encryption.FileEncryptionOpts{
			SourcePath:   "./testdata/txt_to_encrypt.txt",
			OutputPath:   "./testdata/encrypted.txt",
			AESKeyLength: encryption.AES256,
			Key:          key,
			BufferSize:   1024,
		}

		assert.Equal(t, 32, opts.GetKeyLength())
	})

	t.Run("ok get key length", func(t *testing.T) {
		opts := &encryption.FileEncryptionOpts{
			SourcePath:   "./testdata/txt_to_encrypt.txt",
			OutputPath:   "./testdata/encrypted.txt",
			AESKeyLength: encryption.AESKeyLength(10000),
			Key:          key,
			BufferSize:   1024,
		}

		assert.Equal(t, 16, opts.GetKeyLength())
	})

	t.Run("ok get chiper key", func(t *testing.T) {
		opts := &encryption.FileEncryptionOpts{
			SourcePath:   "./testdata/txt_to_encrypt.txt",
			OutputPath:   "./testdata/encrypted.txt",
			AESKeyLength: encryption.AESKeyLength(10000),
			Key:          key,
			BufferSize:   1024,
		}

		res := opts.GetChiperKey()
		assert.Equal(t, 16, len(res))
	})
}

func TestEncryptFile(t *testing.T) {
	key, err := encryption.ReadKeyFromFile("./testdata/test_private.pem")
	assert.NoError(t, err)

	t.Run("file not found / failed open file", func(t *testing.T) {
		opts := &encryption.FileEncryptionOpts{
			SourcePath: "./testdata/nihil.test",
		}

		_, err := encryption.EncryptFile(opts)
		assert.Error(t, err)
	})

	t.Run("failed to create output file", func(t *testing.T) {
		opts := &encryption.FileEncryptionOpts{
			SourcePath:   "./testdata/txt_to_encrypt.txt",
			OutputPath:   "/root/encrypted.txt",
			AESKeyLength: encryption.AES192,
			Key:          key,
			BufferSize:   1024,
		}

		_, err := encryption.EncryptFile(opts)
		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		opts := &encryption.FileEncryptionOpts{
			SourcePath:   "./testdata/txt_to_encrypt.txt",
			OutputPath:   "./testdata/encrypted_text.txt",
			AESKeyLength: encryption.AES192,
			Key:          key,
			BufferSize:   1024,
		}

		_, err := encryption.EncryptFile(opts)
		assert.NoError(t, err)

		err = helper.DeleteFile(opts.OutputPath)
		assert.NoError(t, err)
	})
}

func TestDecryptFile(t *testing.T) {
	key, err := encryption.ReadKeyFromFile("./testdata/test_private.pem")
	assert.NoError(t, err)

	t.Run("invalid file", func(t *testing.T) {
		opts := &encryption.FileEncryptionOpts{
			SourcePath: "./testdata/nihil.test",
		}

		err := encryption.DecryptFile(opts)
		assert.Error(t, err)
	})

	t.Run("failed to create output file", func(t *testing.T) {
		opts := &encryption.FileEncryptionOpts{
			SourcePath:   "./testdata/test_encrypted_text.txt",
			OutputPath:   "/root/encrypted.txt",
			AESKeyLength: encryption.AES192,
			Key:          key,
			BufferSize:   1024,
		}

		err := encryption.DecryptFile(opts)
		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		opts := &encryption.FileEncryptionOpts{
			SourcePath:   "./testdata/test_encrypted_text.txt",
			OutputPath:   "./testdata/decypted.txt",
			AESKeyLength: encryption.AES192,
			Key:          key,
			BufferSize:   1024,
		}

		err := encryption.DecryptFile(opts)
		assert.NoError(t, err)

		assert.FileExists(t, "./testdata/decypted.txt")

		dec, err := os.ReadFile("./testdata/decypted.txt")
		assert.NoError(t, err)

		raw, err := os.ReadFile("./testdata/txt_to_encrypt.txt")
		assert.NoError(t, err)

		assert.Equal(t, dec, raw)

		err = helper.DeleteFile(opts.OutputPath)
		assert.NoError(t, err)
	})
}
