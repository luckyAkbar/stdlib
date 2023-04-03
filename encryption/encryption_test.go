package encryption_test

import (
	"crypto"
	"crypto/rand"
	"testing"

	"github.com/luckyAkbar/stdlib/encryption"
	"github.com/stretchr/testify/assert"
)

func TestEncrypt(t *testing.T) {
	key, err := encryption.GenerateKey(&encryption.KeyGenerationOpts{
		Random: rand.Reader,
		Bits:   2048,
	})
	assert.NoError(t, err)

	t.Run("ok", func(t *testing.T) {
		data := []byte("hello world")
		opts := &encryption.EncryptionOpts{
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
		opts := &encryption.EncryptionOpts{
			Random:    rand.Reader,
			Hash:      crypto.SHA512.New(),
			PublicKey: &key.PublicKey,
		}

		_, err := encryption.Encrypt(data, opts)
		assert.Error(t, err)
	})

	t.Run("ok - encrypt to base64", func(t *testing.T) {
		data := []byte("hello world")
		opts := &encryption.EncryptionOpts{
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
