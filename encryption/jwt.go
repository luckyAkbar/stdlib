package encryption

import (
	"github.com/golang-jwt/jwt/v4"
)

// JWTTokenGenerator is an interface for working with JWT tokens
type JWTTokenGenerator interface {
	// GenerateJWTToken generates a signed JWT token string
	GenerateJWTToken(payload jwt.Claims) (string, error)

	// ValidateJWTToken validates a JWT token string
	// if error is not nil, the token also maybe not nil
	// so if the error is not nil, consider the token as invalid / don't use it
	ValidateJWTToken(token string) (*jwt.Token, error)
}

type jwtToken struct {
	Method     jwt.SigningMethod
	SigningKey []byte
}

// NewJWTTokenHandler creates a new JWTTokenGenerator
func NewJWTTokenHandler(method jwt.SigningMethod, signingKey []byte) JWTTokenGenerator {
	return &jwtToken{
		Method:     method,
		SigningKey: signingKey,
	}
}

func (jtg *jwtToken) GenerateJWTToken(payload jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jtg.Method, payload)

	t, err := token.SignedString(jtg.SigningKey)
	if err != nil {
		return "", err
	}

	return t, nil
}

func (jtg *jwtToken) ValidateJWTToken(token string) (*jwt.Token, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return jtg.SigningKey, nil
	})

	if err != nil {
		return t, err
	}

	return t, nil
}
