// Package helper contains generic helper functions
package helper

import (
	"strings"

	"github.com/google/uuid"
)

// GenerateID generates a random ID using "github.com/google/uuid"
// and removes the "-" from the string
func GenerateID() string {
	id := uuid.New()
	return strings.ReplaceAll(id.String(), "-", "")
}
