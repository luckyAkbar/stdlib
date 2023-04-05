// Package custerr contains custom error types
package custerr

import "fmt"

// ErrChain is a custom error type that can be used to chain errors
// satisfy the error interface
type ErrChain struct {
	Message string
	Cause   error
	Code    int
	Fields  map[string]interface{}
	Type    error
}

// Error returns the error message
func (err ErrChain) Error() string {
	bcoz := ""
	fields := ""
	if err.Cause != nil {
		bcoz = fmt.Sprint(" because {", err.Cause.Error(), "}")
		if len(err.Fields) > 0 {
			fields = fmt.Sprintf(" with Fields {%+v}", err.Fields)
		}
	}
	return fmt.Sprint(err.Message, bcoz, fields)
}
