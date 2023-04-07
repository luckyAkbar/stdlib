package helper_test

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sweet-go/stdlib/helper"
)

func TestHelper_MultipartFileSaver(t *testing.T) {
	file, err := os.Open("./testdata/upload_file.txt")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %s", err)
	}
	defer helper.WrapCloser(file.Close)

	t.Run("failed to open", func(t *testing.T) {
		f := &multipart.FileHeader{}

		err := helper.MultipartFileSaver(f, "")
		assert.Error(t, err)
	})

	t.Run("failed to create output file", func(t *testing.T) {
		// Create a new multipart form and add the file to it
		var buf bytes.Buffer
		writer := multipart.NewWriter(&buf)
		part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
		if err != nil {
			t.Fatalf("Failed to create form file: %s", err)
		}
		if _, err := io.Copy(part, file); err != nil {
			t.Fatalf("Failed to copy file contents to form file: %s", err)
		}
		helper.WrapCloser(writer.Close)
		// Create a new HTTP request with the multipart form data
		req, err := http.NewRequest("POST", "/upload", &buf)
		if err != nil {
			t.Fatalf("Failed to create HTTP request: %s", err)
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())

		// Create a new test server and send the HTTP request
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, fh, err := r.FormFile("file")
			assert.NoError(t, err)

			err = helper.MultipartFileSaver(fh, "/root/output_helper_file_saver.txt")
			assert.Error(t, err)
		})
		handler.ServeHTTP(rr, req)
	})

	t.Run("ok", func(t *testing.T) {
		// Create a new multipart form and add the file to it
		var buf bytes.Buffer
		writer := multipart.NewWriter(&buf)
		part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
		if err != nil {
			t.Fatalf("Failed to create form file: %s", err)
		}
		if _, err := io.Copy(part, file); err != nil {
			t.Fatalf("Failed to copy file contents to form file: %s", err)
		}
		helper.WrapCloser(writer.Close)
		// Create a new HTTP request with the multipart form data
		req, err := http.NewRequest("POST", "/upload", &buf)
		if err != nil {
			t.Fatalf("Failed to create HTTP request: %s", err)
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())

		// Create a new test server and send the HTTP request
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, fh, err := r.FormFile("file")
			assert.NoError(t, err)

			err = helper.MultipartFileSaver(fh, "./output_helper_file_saver.txt")
			assert.NoError(t, err)
			assert.FileExists(t, "./output_helper_file_saver.txt")

			err = helper.DeleteFile("./output_helper_file_saver.txt")
			assert.NoError(t, err)
		})
		handler.ServeHTTP(rr, req)
	})
}
