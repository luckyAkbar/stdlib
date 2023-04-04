package example

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/luckyAkbar/stdlib/encryption"
	"github.com/sirupsen/logrus"
)

// JWTToken is an example of how to use JWT token
func JWTToken() {
	ec := echo.New()

	privatePem := encryption.ParseTestKey(`-----BEGIN RSA TESTING KEY-----
MIIEpQIBAAKCAQEA0Ft5dkFcydEexXT7QD/X+XHRd9PZ4SRLjHP86mJICmu1es7Y
/9+hqEK2sWt2p7KYhYMmcoNqH4FxDAaoURhrHpLEJm/aAzYW8CESSp6E++mqWd/1
lKdGByVEUtq5aP9+T+AHTn+Zuc9vwtTidCtyxoofhL6LsY5n62nsOapDdNOEmu0D
LRvXHNPvK+ctP5DFiY+TOSPoDwTa9RyAeHVLejzl95bbepx5dSoZN69eYvTzNMqn
wEtQnM0OL+WbTXpLC1l8YbccbEVcT7eYbz4FO8BKSMbH7G9Ujb1OUciIL52t7NRj
fwXE+q1h8jGdtoxqMJ9nDJ8wDMpY+WIQ1PXNQwIDAQABAoIBAQC2+J2hi6TAVjR/
kktSEL7I/3rDj/c2D3mIzhK8gbJh9FRalGtbyDdeW7ez9nssqVnnZVTOGzmGkVWF
ChOlo5vuLVSzrGX3i/h4x5IYlTyTPI+sfVBcSsjkXYWyfQF1g2iZHFNOTB/jXJb4
sZpsCfuw/nrPR8XFFxmLUmlv+mVioQaRHbNEEZmcWyHc+6EgHmCYQv65Ymnfj4nD
QzHFkFVh+Whf8rLYU78FFK7i1P4WsDCxVGlh7IXk4vk3uRLNxti8ytMLGzgJYtdf
6SMsR0fqYOicJtgtiNVOm/T1gNP/Ix4/BD/7+voMdBfeEpl5c7CpbhqCDDfj7RbQ
jS3n2QeBAoGBANbeDcefCirIW+5IP1FbZ68WBnFPW6SUgdpPbBtgcPInWJ/9kia2
fJgWqW49yCfqkOOy6dHm57o1ECLAUPN+j9tjb4SoMK21epn0JYHsYiXnPkLetRN8
zBayyS+ULzvCigUppqBoB0sw2+7TYNY1RbbsjRIwTcpne0io4Z0nXYRlAoGBAPg+
Yhh1xz099TM96TOpzY8yrvbMz5Rj+uIaR5wR0NgRa5fsYCXOjqj2gZNQbdz1BJaV
ABpJdManCkpUYiYa9DqL7Zkc4wLe0AW+bJFBjJHjqd8Aqw7Nk/+ehjeSlg25n+Bo
fx1C//GI/XwmPevRyFPpDNV3AzUtk/rnFKKzC0yHAoGAB7HeFncImy2ftTHbKqO2
W9vTET3BT2yOFe5gNb7HbLSiBODE2iQQ5DVzjeIih+NrmuvuWbkGNXHvCP+QJpgy
uK2f8cVAMQhdwqOusC9x+F+GqEhnfbIrcOioMc8BvgcigDrUn8v57uRqC+x//Eve
GkXwa2VVc9ku3hRGOCWPwM0CgYEAhbT8Gxac+NR9VFs9VzFXYZC4AoBwMgnj4JKt
DVffN/GyFQMhClwGJOWZByKj+gYSsZSRmJcGCdWAymZG8yVDdKFXmUeg0jP2sZFO
YrJ+pzmLjmyKtg9ubpkQy6/tmHjprvI5vSYQOyVA+vSSF4lHsEJvQi63EJZ7BQIf
8D4lkNMCgYEAnCwqm+vzBF3F372n7hovr53OUUopggHkaSCG8y3O3vXcoXwDdgZL
8tUELgfxf3eAIl1StvBzVjJJbe4asl6muU/upOPrHsnIrl8lfiAthmQ1DD+zlNQQ
ugO2voQz8H++ZICpQiCns4CdLtPC1fAq2LOzinOq45Hb3XFWT7c1N5k=
-----END RSA TESTING KEY-----
`)

	jwtgen := encryption.NewJWTTokenHandler(jwt.SigningMethodHS256, []byte(privatePem))

	ec.GET("/login", func(c echo.Context) error {
		type custClaim struct {
			jwt.RegisteredClaims
			UserID string `json:"user_id"`
		}

		token, err := jwtgen.GenerateJWTToken(&custClaim{
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "example",
				Subject:   "example only",
				Audience:  []string{"example"},
				ExpiresAt: jwt.NewNumericDate(jwt.TimeFunc().Add(5 * time.Minute)),
				NotBefore: jwt.NewNumericDate(jwt.TimeFunc().Add(1 * time.Minute)),
				IssuedAt:  jwt.NewNumericDate(jwt.TimeFunc()),
				ID:        "example-id",
			},
			UserID: "test user id",
		})

		if err != nil {
			logrus.Info(err)
			return c.NoContent(http.StatusInternalServerError)
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"token": token,
		})
	})

	g := ec.Group("/protected")
	g.Use(jwtgen.BuildEchoJWTMiddleware())

	g.GET("/p", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "protected",
		})
	})

	ec.Logger.Fatal(ec.Start(":8080"))
}
