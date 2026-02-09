// Package service is the shared backage by Gophermar services
package service

import (
	"context"

	"github.com/google/uuid"
	genBalance "github.com/oleshko-g/gophermart/internal/gen/balance"
	genUser "github.com/oleshko-g/gophermart/internal/gen/user"
)

// Service is the struct to set up [balance.Service] implementation
type Service struct {
	User genUser.Service
	Balance
}

// Auther is the interface to be met by authenticating endpoints of Gophermart services
type Auther interface {
	genBalance.Auther
	UserIDFromContext(context.Context) (uuid.UUID, error)
}

// Balance is the interface to be met by a [balance.Service] implementation
type Balance interface {
	genBalance.Service
	ProcessAccruals(ctx context.Context) error
}
