package balance

import (
	"github.com/oleshko-g/oggophermart/internal/service/errors"
)

var (
	// ErrOwnerMismatch is the error value which is used to map to the 409 Conflict HTTP Status code
	ErrOwnerMismatch = errors.New("The order belongs to another user")

	// ErrInvalidOrderNumber is the error value which is used to map to the 422 Unprocessable Entity HTTP Status code
	ErrInvalidOrderNumber = errors.New("Invalid order number")

	// ErrAccrualNotStarted is the error returned when new accrual order arrives while accrual proccessing hasn't started yet
	ErrProcessAccrualsNotStarted = errors.New("Accrual processing hasn't started yet")

	// ErrFailedToGetOrderAccrual is the error returned when accrual order while accrual proccessing hasn't started yet
	ErrFailedToGetOrderAccrual = errors.New("Failed to get order accrual")

	// ErrUnknownAccrualOrderStatus is the error return when accrual order status isn't known
	ErrUnknownAccrualOrderStatus = errors.New("Accrual order status isn't known")

	// ErrNegativeBalance is the error returned when a user balance have become negative after a withdrawal
	ErrInsufficientFunds = errors.New("Insufficient funds")

	// ErrCurrentBalanceChanged is the error returned when a user balance has been changed by other transaction since current transaction changed it
	ErrCurrentBalanceChanged = errors.New("Current balance has been changed by other transaction")
)
