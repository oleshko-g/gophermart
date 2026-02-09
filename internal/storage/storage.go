// Package storage declares the type to imported by the service architecture layer
package storage

import (
	"context"
	"time"

	"github.com/google/uuid"
	genDBSQL "github.com/oleshko-g/gophermart/internal/gen/storage/db/sql"
)

// Storage is the struct to hold an implementation of gophermart storage
type Storage struct {
	Transacter //interface
	User       // interface
	Balance    // interface
}

// Order is an alias for [genDBSQL.Order]
type Order = genDBSQL.Order

// User declares the storage interface for the user service
type User interface {
	RetrieveUser(ctx context.Context, login string) (userID uuid.UUID, err error)
	RetreiveUserPassword(ctx context.Context, login string) (hashedPassword string, err error)
	StoreUser(ctx context.Context, login, hashedPassword string) error
}

// Balance declares the storage interfce for the balance service
type Balance interface {
	Transacter
	StoreOrder(ctx context.Context, userID uuid.UUID, orderNumber, status string, createdAt time.Time) (orderID uuid.UUID, err error)
	RetreiveOrderUser(ctx context.Context, orderNumber string) (userID uuid.UUID, err error)
	RetrieaveUserOrders(ctx context.Context, userID uuid.UUID) ([]genDBSQL.SelectOrdersByUserIDRow, error)
	Retrieve(ctx context.Context, userID uuid.UUID) (genDBSQL.SelectBalanceByUserIDRow, error)
	RetrieveOrderIDsForAccrual(ctx context.Context) ([]uuid.UUID, error)
	RetrieveOrderForAccrual(ctx context.Context, orderID uuid.UUID) (Order, error)
	UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status string) error
	StoreUserAccrual(ctx context.Context, userID uuid.UUID, orderID uuid.UUID, amount int32) error
	StoreUserWithdrawal(ctx context.Context, userID uuid.UUID, orderID uuid.UUID, amount int32) (withdrawalID uuid.UUID, err error)
	RetrieveUserWithdrawals(ctx context.Context, userID uuid.UUID) (withdrawals []genDBSQL.SelectBalanceOrderTransactionAmountByUserIDAndKindRow, err error)
}

// Transaction is the interface that must be met by a storage that supports storage transactions
type Transaction interface {
	Commit() error
	Rollback() error
}

// Tx is the type returned by [Transacter.BeginTx]
type Tx struct {
	Tx Transaction
	Balance
}

// Transacter is the interface that must be met by a storage that supported storage transactions
type Transacter interface {
	BeginTx(context.Context) (*Tx, error)
}
