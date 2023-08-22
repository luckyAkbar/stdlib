// Package helper contains generic helper functions
package helper

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

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

// LogIfError log error if error is not nil
func LogIfError(err error) {
	if err != nil {
		logrus.Error(err)
	}
}

// GenerateUniqueName generate unique name using GenerateID and time.Now().Format(time.RFC3339)
func GenerateUniqueName() string {
	return GenerateID() + time.Now().Format(time.RFC3339)
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

// ReadFileMetadata read file metadata from multipart file header
func ReadFileMetadata(file *multipart.FileHeader) (*FileMetadata, error) {
	f, err := file.Open()
	if err != nil {
		return nil, err
	}

	buffer := make([]byte, 512)
	_, err = f.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return &FileMetadata{
		Name:        file.Filename,
		Size:        file.Size,
		ContentType: http.DetectContentType(buffer),
	}, nil
}

// DeleteFile wrapper for os.Remove. will delete file from given path
func DeleteFile(path string) error {
	return os.Remove(path)
}
