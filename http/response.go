// Package http is a package containing all functionalities related to HTTP
package http

import (
	"encoding/json"

	"github.com/luckyAkbar/stdlib/encryption"
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

// APIResponseGenerator is a struct containing functionalities to generate API response with signature
type APIResponseGenerator struct {
	defaultEncryptionOpts *encryption.SignOpts
}

// NewStandardAPIResponseGenerator is a constructor for APIResponseGenerator
func NewStandardAPIResponseGenerator(defaultEncryptionOpts *encryption.SignOpts) *APIResponseGenerator {
	return &APIResponseGenerator{
		defaultEncryptionOpts: defaultEncryptionOpts,
	}
}

// GenerateAPIResponse is a function to generate API response with signature.
// If opts is nil, it will use defaultEncryptionOpts
func (arg *APIResponseGenerator) GenerateAPIResponse(response *StandardResponse, opts *encryption.SignOpts) (*APIResponse, error) {
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
