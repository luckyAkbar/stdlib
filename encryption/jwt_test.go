package encryption_test

import (
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/luckyAkbar/stdlib/encryption"
	"github.com/stretchr/testify/assert"
)

func TestGenerateJWTToken(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		jwtgen := encryption.NewJWTTokenHandler(jwt.SigningMethodHS256, []byte("secret"))

		token, err := jwtgen.GenerateJWTToken(jwt.RegisteredClaims{
			Issuer: "test",
		})

		assert.NoError(t, err)

		_, err = jwtgen.ValidateJWTToken(token)

		assert.NoError(t, err)
	})

	t.Run("bad signing key", func(t *testing.T) {
		key, err := encryption.GenerateKey(nil)
		assert.NoError(t, err)

		jwtgen := encryption.NewJWTTokenHandler(jwt.SigningMethodRS256, []byte(key.D.String()))

		_, err = jwtgen.GenerateJWTToken(jwt.RegisteredClaims{
			Issuer: "test",
		})

		assert.Error(t, err)
	})

	t.Run("bad token", func(t *testing.T) {
		jwtgen := encryption.NewJWTTokenHandler(jwt.SigningMethodHS256, []byte("secret"))

		_, err := jwtgen.ValidateJWTToken("bad token")

		assert.Error(t, err)
	})
}
