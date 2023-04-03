// Package example containing sample code for using this library
package example

import (
	"crypto"
	"crypto/rand"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/luckyAkbar/stdlib/encryption"
	stdlib_http "github.com/luckyAkbar/stdlib/http"
)

// APIResponse is example how to run API and see the API response with signature
func APIResponse() {
	ec := echo.New()
	privatePem := `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAwiSX09qKwzg+eunMwn4AulMCHc2z77jy2Mx0Ehc4x014l0Oz
W5+V5rYadipEM2gLLCdb2nE9hY0+0zC3GEoxoT5ksLdDw6kSOiI5iYQsKkULG9aT
2g4Bw14VPLiEt2jpzU3D80XueVsQD6OvfG9LErMEkIggTtfJWUx37SUFgddq/XOD
x2VVll5p4Aoe4I7z1Z2AXzXI0CeCG7Q8taeSwKVlqzKwjYbCN+64JZgmCgJ/JIxI
PLfrdCszZ8MhkT4fKMkFZSKIIOA3YsTkDThrCmGKiwYbHqKZkBPQILHXsmFchTD1
jCFGxNLGnuFKTdL1HdRtGDF0bUerpqFrLPf8jQIDAQABAoIBAB8ZQYDcJxIydj+2
J+iXyoIOPe6MPuCNnckApy8mrR+v1ztTyu1IWPjF/uMizh318qQ2Ac6yCQrVn1Sx
HwMzD1Qm7HYqRR6YfKT3SoQuuexjdu4Up0Zsq/ehoIFDhW7jzV/nrxXPA+5ImgAH
Vlr2cO4j4v1L8PDwO/6j8yn3njQ4B5CLK29pEbKSP+plczEHq9wbBQeeanvy6yYh
8fSeQOzQu5SgeQDcVKAZG+6TfuUVAL+G2E8LiPhOeOg/T7coFAuChLot0gukQwo0
sXyEI75IflsXiTjgWpq/0lYAUxpHVJ5JBy3CXqYS+95J9b2/T9k9cZ3xHQWqi/iT
qTVqqTkCgYEA/hXaJOP2ho45UTbkluSvTKk3kkIGn0T5a+D8X1MGuQ1tn00jaGnX
W2y3yRxuNMg4ZbdePuhnCw1FXd59Cy+FydUkyU6TBsaEYUvLnOagPb7OQOU5+N6a
AVN2AIpVNzFIlenvZwAPHJitFFX0dDXxy57/DC8HxKkEviZVVdK54h8CgYEAw5sb
o5fGw9K94nXOVp2gbIHS5/QI7FOxZqCfkxfO5ssfxRDzfcFUdK7ijw/PiJR1iMIZ
s+8Fn0vX55zf3Slk9Du629h2I4fB2ViT12Lu+7u+Sc3XV+/pr7e0XZCKvNNc/B4q
pb44AFDhhATLxMlhiailzkVSlQVhqZWSBsrbw9MCgYEAzKphe1G9JImvlcG3w+wV
YJT11HQmzWhL2R/zaf2A7tLoOGd0XAjVlikuqWqjQxT8iMJ5wgaF6hsYgxJSew4e
oIN2DEmkmNKTk6PwMUR8UwA9N3ztg5AbUXIfHTRQjBLAuzEizD757Tj2Qeky7eD+
EdzS6MeBZGIZFene1zDU1lUCgYBz1/+VckTgIoYcgUJzX6TrvjNG9er281YD/qqi
9Z2uZ6voDPL3jjDTbeN1cJqrO6kkFjgcrTk6LzOt0uVt2J8WWe1/WAIXZsYyT1g3
XjtE0NqQYRzg0pAmZfim1PyledP+6Gq/gBkwbrYwdpqrb8yZN00DDWEsKmS9h3xV
E3z1ywKBgHSt/WzI0sp1lBmm2N/S3QAmvv2XaGHO0Zlr8fvcgwHuu7X1NBLEuUri
euJMSq2sQEtD49W2+9DKcOvI+qnlyumeOryCY2NTRTeFznUBwPpnx0+hvC3xF9j+
SGUVK+7JIUp8ee5gCqsJrzK82j00IH6kk6zrX6zvQN3Zd5O4ImZd
-----END RSA PRIVATE KEY-----
`
	key, err := encryption.ReadKey([]byte(privatePem))

	if err != nil {
		panic(err)
	}

	apiResponseGen := stdlib_http.NewStandardAPIResponseGenerator(&encryption.SignOpts{
		Random:  rand.Reader,
		PrivKey: key,
		Alg:     crypto.SHA512,
		PSSOpts: nil,
	})

	type data struct {
		Message string `json:"message"`
	}

	response := stdlib_http.StandardResponse{
		Success: true,
		Message: "Success",
		Status:  http.StatusOK,
		Data:    data{Message: "Hello World"},
	}

	ec.GET("/test", func(c echo.Context) error {
		resp, err := apiResponseGen.GenerateAPIResponse(&response, nil)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		return c.JSON(http.StatusOK, resp)
	})

	ec.Logger.Fatal(ec.Start(":8080"))
}
