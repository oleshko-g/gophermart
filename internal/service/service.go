// Package service is
package service

import (
	"context"

	"github.com/google/uuid"
	genBalance "github.com/oleshko-g/gophermart/internal/gen/balance"
	genUser "github.com/oleshko-g/gophermart/internal/gen/user"
)

type Service struct {
	User genUser.Service
	Balance
}

type Auther interface {
	genBalance.Auther
	UserIDFromContext(context.Context) (uuid.UUID, error)
}

type Balance interface {
	genBalance.Service
	ProcessAccruals(ctx context.Context) error
}
