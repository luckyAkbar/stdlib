// Package helper contains generic helper functions
package helper

import (
	"strings"

	"github.com/google/uuid"
	"github.com/kumparan/go-utils"
	"github.com/sirupsen/logrus"
)

// GenerateID generates a random ID using "github.com/google/uuid"
// and removes the "-" from the string
func GenerateID() string {
	id := uuid.New()
	return strings.ReplaceAll(id.String(), "-", "")
}

// Dump to json using json marshal. wrapper for "github.com/kumparan/go-utils".Dump func
func Dump(i interface{}) string {
	return utils.Dump(i)
}

// WrapCloser wrap closer. If closer return error, log the error
func WrapCloser(closeFn func() error) {
	if err := closeFn(); err != nil {
		logrus.Error(err)
	}
}
