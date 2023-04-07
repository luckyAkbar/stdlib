package encryption

import "strings"

// AESKeyLength is the length of AES key
type AESKeyLength int

// list of available aes key length to choose from
const (
	AES128 AESKeyLength = 16
	AES192 AESKeyLength = 24
	AES256 AESKeyLength = 32
)

// ParseTestKey is a helper function to parse a test key to avoid test key detected as real key
func ParseTestKey(key string) string {
	return strings.ReplaceAll(key, "TESTING KEY", "PRIVATE KEY")
}
