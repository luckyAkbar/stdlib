package http

import (
	"encoding/json"

	"github.com/luckyAkbar/stdlib/encryption"
)

type StandardResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Status  int    `json:"status"`
	Data    any    `json:"data"`
	Error   string `json:"error,omitempty"`
}

type APIResponse struct {
	Response  any    `json:"response"`
	Signature string `json:"signature"`
}

type APIResponseGenerator struct {
	defaultEncryptionOpts *encryption.SignOpts
}

func NewStandardAPIResponseGenerator(defaultEncryptionOpts *encryption.SignOpts) *APIResponseGenerator {
	return &APIResponseGenerator{
		defaultEncryptionOpts: defaultEncryptionOpts,
	}
}

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
