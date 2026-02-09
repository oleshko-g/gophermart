// Package transport defines [Accrual] interface and type for working with it
package transport

import (
	"context"
)

// Accrual is the interface to accrual system
//
//go:generate moq -pkg moqAccrual -out ../gen/transport/moq/accrual/accrual.go . Accrual
type Accrual interface {
	// FetchOrderAccrual method
	FetchOrderAccrual(ctx context.Context, payload FetchOrderAccrualPayload) (*FetchOrderAccrualResult, error)
}

// FetchOrderAccrualPayload is the payload type of the accrual system of [FetchOrderAccrual] method
// FetchOrderAccrual method
type FetchOrderAccrualPayload struct {
	Number string
}

// FetchOrderAccrualResult is the result type of the accrual system of [FetchOrderAccrual] method
type FetchOrderAccrualResult struct {
	Order   string
	Status  OrderAccrualStatus
	Accrual *float64
}

// AccrualError is the type to serialize the HTTP response of the Accrual service
type AccrualError struct {
	// identifier to map an error to HTTP status codes
	RetryAfter int
	Message    string
}

// Error returns an error Message.
func (e *AccrualError) Error() string {
	return e.Message
}

// OrderAccrualStatus is the type to store known accrual statuses
type OrderAccrualStatus string

// Known accrual statuses
const (
	OrderAccrualStatusRegistered = "REGISTERED"
	OrderAccrualStatusProcessing = "PROCESSING"
	OrderAccrualStatusProcessed  = "PROCESSED"
	OrderAccrualStatusInvalid    = "INVALID"
)
