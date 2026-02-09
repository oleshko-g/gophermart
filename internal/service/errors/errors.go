// Package errors is the package to errors shared by all gophermart services
package errors //revive:disable-line:var-naming

import (
	genSvc "github.com/oleshko-g/gophermart/internal/gen/service"
)

type svcError = genSvc.GophermartError

// New returns a pointer to a new instance of gophermart service error
func New(name string) error {
	return &svcError{Name: name}
}

var (
	// ErrUserIsNotAuthenticated is the error value which is used to map to the 401 Unauthorized HTTP Status code
	ErrUserIsNotAuthenticated = New("User is not authenticated")
	// ErrInternalServiceError is the error value which is used to map to the 500 Internal Server Error HTTP Status code
	ErrInternalServiceError = New("Internal service error")
)
