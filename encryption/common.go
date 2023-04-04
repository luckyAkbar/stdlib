package encryption

import "strings"

// ParseTestKey is a helper function to parse a test key to avoid test key detected as real key
func ParseTestKey(key string) string {
	return strings.ReplaceAll(key, "TESTING KEY", "PRIVATE KEY")
}
