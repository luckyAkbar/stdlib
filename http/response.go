// Package http is a package containing all functionalities related to HTTP
package http

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sweet-go/stdlib/encryption"
	custerr "github.com/sweet-go/stdlib/error"
)

// StandardResponse is a standard response for all API
type StandardResponse struct {
	Success bool   `json:"success,omitempty"`
	Message string `json:"message,omitempty"`
	Status  int    `json:"status,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

// APIResponse is a standard response for all API
type APIResponse struct {
	Response  any    `json:"response"`
	Signature string `json:"signature"`
}

// APIResponseGenerator is an interface containing functionalities to generate standard API response
type APIResponseGenerator interface {
	GenerateAPIResponse(response *StandardResponse, opts *encryption.SignOpts) (*APIResponse, error)

	// GenerateEchoAPIResponse is a function to generate API response with signature and send it to echo context also act as a wrapper for GenerateAPIResponse
	// will use http status code the same as response.Status
	GenerateEchoAPIResponse(c echo.Context, response *StandardResponse, opts *encryption.SignOpts) error
}

// APIResponseGenerator is a struct containing functionalities to generate API response with signature
type apiResponseGenerator struct {
	defaultEncryptionOpts *encryption.SignOpts
}

// NewStandardAPIResponseGenerator is a constructor for APIResponseGenerator
func NewStandardAPIResponseGenerator(defaultEncryptionOpts *encryption.SignOpts) APIResponseGenerator {
	return &apiResponseGenerator{
		defaultEncryptionOpts: defaultEncryptionOpts,
	}
}

// GenerateAPIResponse is a function to generate API response with signature.
// If opts is nil, it will use defaultEncryptionOpts
func (arg *apiResponseGenerator) GenerateAPIResponse(response *StandardResponse, opts *encryption.SignOpts) (*APIResponse, error) {
	if opts == nil {
		opts = arg.defaultEncryptionOpts
	}

	respBytes, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}

	signature, err := encryption.SignToBase64(respBytes, opts)
	if err != nil {
		return nil, err
	}

	return &APIResponse{
		Response:  response,
		Signature: signature,
	}, nil
}

func (arg *apiResponseGenerator) GenerateEchoAPIResponse(c echo.Context, response *StandardResponse, opts *encryption.SignOpts) error {
	resp, err := arg.GenerateAPIResponse(response, opts)
	if err != nil {
		return &custerr.ErrChain{
			Message: "Failed to generate API response",
			Cause:   err,
			Code:    http.StatusInternalServerError,
		}
	}

	return c.JSON(response.Status, resp)
}
