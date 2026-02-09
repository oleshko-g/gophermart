// Package errors is the internal storage package
package errors //revive:disable:var-naming

import (
	"errors"
)

var (
	// ErrUnsupportedDataSource is returned when data source is not supported by the storage implementation
	ErrUnsupportedDataSource = errors.New("unsupported data source")

	// ErrMissingDatabaseName is returned when data source lacks database name
	ErrMissingDatabaseName = errors.New("missing database name")

	// ErrNoAffect is returned when storage operation didn't do anything
	ErrNoAffect = errors.New("no affect")

	// ErrAlreadyExists is returned when a storage record alrady exists
	ErrAlreadyExists = errors.New("already exists")

	// ErrNotFound is returned when a storage record is not found
	ErrNotFound = errors.New("not found")
)
