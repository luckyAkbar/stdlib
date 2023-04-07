// Package helper contains generic helper functions
package helper

import (
	"io"
	"mime/multipart"
	"os"
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

// MultipartFileSaver save multipart file to given path
func MultipartFileSaver(file *multipart.FileHeader, path string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer WrapCloser(src.Close)

	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer WrapCloser(dst.Close)

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return nil
}

// DeleteFile wrapper for os.Remove. will delete file from given path
func DeleteFile(path string) error {
	return os.Remove(path)
}
