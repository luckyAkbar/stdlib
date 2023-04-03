package encryption_test

import (
	"crypto/rand"
	"os"
	"testing"

	"github.com/luckyAkbar/stdlib/encryption"
	"github.com/stretchr/testify/assert"
)

func TestGenerateKey(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		cfg := &encryption.KeyGenerationOpts{
			Random:          rand.Reader,
			Bits:            2048,
			PEMFormat:       true,
			GOBFormat:       true,
			PublicFilename:  "_test_public",
			PrivateFilename: "_test_private",
		}
		_, err := encryption.GenerateKey(cfg)
		assert.NoError(t, err)

		assert.FileExists(t, "_test_public.pem")
		assert.FileExists(t, "_test_private.pem")
		assert.FileExists(t, "_test_public.key")
		assert.FileExists(t, "_test_private.key")

		err = os.Remove("./_test_public.pem")
		assert.NoError(t, err)

		err = os.Remove("./_test_private.pem")
		assert.NoError(t, err)

		err = os.Remove("./_test_public.key")
		assert.NoError(t, err)

		err = os.Remove("./_test_private.key")
		assert.NoError(t, err)
	})

	t.Run("random is nil", func(t *testing.T) {
		cfg := &encryption.KeyGenerationOpts{
			Random: nil,
		}

		_, err := encryption.GenerateKey(cfg)
		assert.Error(t, err)
	})
}
