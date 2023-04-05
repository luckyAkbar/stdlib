package encryption_test

import (
	"crypto/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sweet-go/stdlib/encryption"
)

func TestReadKeyFromFile(t *testing.T) {
	t.Run("file not exists", func(t *testing.T) {
		filename := "./_test_notexits.pem"

		_, err := encryption.ReadKeyFromFile(filename)

		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		_, err := encryption.GenerateKey(&encryption.KeyGenerationOpts{
			Random:          rand.Reader,
			Bits:            2048,
			PEMFormat:       true,
			PublicFilename:  "_test_public_TestReadKeyFromFile",
			PrivateFilename: "_test_private_TestReadKeyFromFile",
		})
		assert.NoError(t, err)

		_, err = encryption.ReadKeyFromFile("_test_private_TestReadKeyFromFile.pem")
		assert.NoError(t, err)

		assert.FileExists(t, "_test_public_TestReadKeyFromFile.pem")
		assert.FileExists(t, "_test_private_TestReadKeyFromFile.pem")

		err = os.Remove("./_test_public_TestReadKeyFromFile.pem")
		assert.NoError(t, err)

		err = os.Remove("./_test_private_TestReadKeyFromFile.pem")
		assert.NoError(t, err)
	})
}

func TestReadKey(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		key := encryption.ParseTestKey(`-----BEGIN RSA TESTING KEY-----
MIIEpAIBAAKCAQEA37/7OzLTSBubEeRWSKSRHIJwk+uAnGFFJenKo9TXPh1oa/20
l92hsVeWk76+oo6jRqHskiOUoYYlxiofriomZlM8OyDrDNndaE+havKvvsfksLzq
sAgM/JEfXOpzrAooI0I07EKNI7smhRRrivOY0jLETBfdZ0r729tPvahgKzf4kpiU
Hd9/jK9qFykHt/EhhDq9mPS24S7CcyMm9HysV9b1pViCQXJFA9OrOewVP1wWXRVE
afXngJE5nVpso8B+HfUbsVMPOP9bBps0FJQHBrxvTEobj0ulQbaz8sAWGa7JHXf9
sXp5vUUtLPjxcDsznCbqvkSeIOkh2MBRJSzDNwIDAQABAoIBAHY1t0VPVNCDxSlu
uScnyoKFZ3S+tvPnb+DX43cqu4zVfJWRNBgHv6Ux4Rutaon3Ucu/QHz0z1GGze4j
0xjwq9jjoK6cdZIUiCTT7TPTg4YHlYrKRDM8DaBiC2/LbdE2jH4UPGGVx3tZJMCq
SSbgC50BtTN+aDpqIyXEeBx7GFO8ATl57qZcT10/uAYR8P/w6T50SgV92TjTwrXY
/upoXQ5PGzuBOS0VkEcXquAhYsq0u/b7lzqlpEJMxQjStvqeZOqQEt4GdHJIy3/P
7EfaInNV22Z7bPEjdh7ZfPPP9d8bbgEldwqIdrPoCsghxmniehcZcLx0ltVOSj6z
QJKTWgECgYEA5V3aD+1x1VINRlZJA7ijehjrU7F4Yi7f4PIeh7XnfBBcQiJVn0u7
ti2IwKymIVG6XXRXIlH5f+MIYEXzgE3bxIgNuyCpeh2qchauKR2eh2T3w/zI+/Ce
yuFSrizaNXZ4FkzYEm2iUwiGLFVoumhz7EP7cNPOFSOhEMsEVQ6l2NsCgYEA+bsr
Dh/JEW8L+orbWYOMSOcR1dbdlR7Pn6iP00Vj67wIWnjdHUtAwE+b79FStSPgIpUM
lDTfyMx96UU+s+0oodwEvAZNLBxMaKMQrgrJLatKAINjKx7KX80cXTdetTqCUC2z
YLyDcI0+TZ36Go1kHUb9hKM1F1lH1VWty043j9UCgYEAl4Yjy8fiHrng+SmBfMra
fIu/0v939uzei72Hu8HJFiW8vRfvlpeyf0yffiHQckyKoLh947dh60FxxCASGB3X
ZIM5Bvkx3PGCK3KeRZ1CoFFsePYjVIUGciLeux/4W79S3/COAcaZqN8FvH4D/LmK
c3gJwOS7zS1Hd0+XIhXWLGcCgYBndbpVpKd5WIce6g4L3KruvQQvkk/Earpbi8ri
HTpTPFg9mxsH+tg9k/2nchIQx2chDJzkfa9EkiuLy8s5YYRW4j734qhwIN0q8HuF
jyRfjjofUk9wWtY+sEwS9lB/RlkcfIJ3DkJqC6oHH+6wt2kFlBaNr8vb+3n+EPvq
YWI1bQKBgQDRwc3Jrm/+zjfR3BypS9GW/AV4PG5+Ua+djN1rsSkXKthcnA68bXe6
YXY0IRpoR2MUhC13PqruU9u9WBp0fcYtU/yWJJkCRaYSsBEfDZe+553O/Hr+SQjm
Bfg6AwEBACcmHq/3z/EDdeol5D1sAg223xEiAD+ABceWHpzoaCNyRg==
-----END RSA TESTING KEY-----
`)

		_, err := encryption.ReadKey([]byte(key))

		assert.NoError(t, err)
	})

	t.Run("not ok", func(t *testing.T) {
		// pkcs8 format
		key := encryption.ParseTestKey(`-----BEGIN TESTING KEY-----
MIIBVgIBADANBgkqhkiG9w0BAQEFAASCAUAwggE8AgEAAkEAq7BFUpkGp3+LQmlQ
Yx2eqzDV+xeG8kx/sQFV18S5JhzGeIJNA72wSeukEPojtqUyX2J0CciPBh7eqclQ
2zpAswIDAQABAkAgisq4+zRdrzkwH1ITV1vpytnkO/NiHcnePQiOW0VUybPyHoGM
/jf75C5xET7ZQpBe5kx5VHsPZj0CBb3b+wSRAiEA2mPWCBytosIU/ODRfq6EiV04
lt6waE7I2uSPqIC20LcCIQDJQYIHQII+3YaPqyhGgqMexuuuGx+lDKD6/Fu/JwPb
5QIhAKthiYcYKlL9h8bjDsQhZDUACPasjzdsDEdq8inDyLOFAiEAmCr/tZwA3qeA
ZoBzI10DGPIuoKXBd3nk/eBxPkaxlEECIQCNymjsoI7GldtujVnr1qT+3yedLfHK
srDVjIT3LsvTqw==
-----END TESTING KEY-----
`)

		_, err := encryption.ReadKey([]byte(key))

		assert.Error(t, err)
	})
}

func TestReadPublicKey(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		pub := `-----BEGIN PUBLIC KEY-----
MIIBCgKCAQEAuRozdJRuXNPUocFebG/lTpWvXjk7ciCBomBPybVwPMSbBxonBRbu
lgZ9OhwYsqp99BAyaeTikTHkHAW+oP7LZF5AGN0NPZ66Py9G6zhAmPnoBUcp/CV0
64lJj7ykd6b6EmZyKw3X+uMwwr46bFU9m/Nx29L3yUIZlrnbnoVjzJbSmZXRJVTt
2aRSlv7+aFr20HbIyMtA/+QTm4T+KZuKBl1BBe0edCkNDrZQVqP1w/kjIfOt0s0t
hxhnxBsM64h/83CiWb3zqLhDpOdpe88Xw6VShLEc+cts5b6VOjwgK9sn4rOmP13g
JdJGprKgw66DT6swVH1+x1zgVbpgGdKc+QIDAQAB
-----END PUBLIC KEY-----
`

		_, err := encryption.ReadPublicKey([]byte(pub))
		assert.NoError(t, err)
	})

	t.Run("not ok", func(t *testing.T) {
		pub := `-----BEGIN RSA PUBLIC KEY-----
MIIBCgKCAQEAuRozdJRuXNPUocFebG/lTpWvXjk7ciCBomBPybVwPMSbBxonBRbu
lgZ9OhwYsqp99BAyaeTikTHkHAW+oP7LZF5AGN0NPZ66Py9G6zhAmPnoBUcp/CV0
64lJj7ykd6b6EmZyKw3X+uMwwr46bFU9m/Nx29L3yUIZlrnbnoVjzJbSmZXRJVTt
2aRSlv7+aFr20HbIyMtA/+QTm4T+KZuKBl1BBe0edCkNDrZQVqP1w/kjIfOt0s0t
hxhnxBsM64h/83CiWb3zqLhDpOdpe88Xw6VShLEc+cts5b6VOjwgK9sn4rOmP13g
JdJGprKgw66DT6swVH1+x1zgVbpgGdKc+QIDAQAB
-----END RSA PUBLIC KEY-----
`

		_, err := encryption.ReadPublicKey([]byte(pub))
		assert.Error(t, err)
	})

	t.Run("not ok", func(t *testing.T) {
		pub := `-----BEGIN PUBLIC KEY-----
	MIIBCgKCAQEAuRozdJRuXNPUocFebG/lTpWvXjk7ciCBomBPybVwPMSbBxonBRbu
	lgZ9OhwYsqp99BAyaeTikTHkHAW+oP7LZF5AGN0NPZ66Py9G6zhAmPnoBUcp/CV0
	64lJj7ykd6b6EmZyKw3X+uMwwr46bFU9m/Nx29L3yUIZlrnbnoVjzJbSmZXRJVTt
	2aRSlv7+aFr20HbIyMtA/+QTm4T+KZu1BBe0edCkNDrZQVqP1w/kjIfOt0s0t
	hxhnxBsM64h/83CiWb3zqLhDpOdpe88Xw6VShLEc+cts5b6VOjwgK9sn4rOmP13g
	JdJGprKgw66DT6swVH1+x1zgVbpgGdKc+QIDAQAB
	-----END PUBLIC KEY-----
	`

		_, err := encryption.ReadPublicKey([]byte(pub))
		assert.Error(t, err)
	})

	t.Run("read from file", func(t *testing.T) {
		key, err := encryption.GenerateKey(&encryption.KeyGenerationOpts{
			Random:          rand.Reader,
			Bits:            2048,
			PEMFormat:       true,
			PublicFilename:  "__testing_readfile_public",
			PrivateFilename: "__testing_readfile_private",
		})
		assert.NoError(t, err)
		assert.FileExists(t, "./__testing_readfile_public.pem")
		assert.FileExists(t, "./__testing_readfile_private.pem")

		pub, err := encryption.ReadPublicKeyFromFile("./__testing_readfile_public.pem")
		assert.NoError(t, err)

		assert.Equal(t, key.PublicKey, *pub.PublicKey)

		err = os.Remove("./__testing_readfile_public.pem")
		assert.NoError(t, err)

		err = os.Remove("./__testing_readfile_private.pem")
		assert.NoError(t, err)
	})

	t.Run("read from file failed", func(t *testing.T) {
		_, err := encryption.ReadPublicKeyFromFile("./imaginary_key_never_exists.pem")
		assert.Error(t, err)
	})
}
