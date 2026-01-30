// Package sql is the impementation of
package sql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	pgDriver "github.com/lib/pq" // revive:disable-line:blank-imports registers the postgres driver
	genDBSQL "github.com/oleshko-g/oggophermart/internal/gen/storage/db/sql"
	"github.com/oleshko-g/oggophermart/internal/storage"
	"github.com/oleshko-g/oggophermart/internal/storage/db"
	"github.com/oleshko-g/oggophermart/internal/storage/db/sql/schema"
	storageErrors "github.com/oleshko-g/oggophermart/internal/storage/errors"
)

// New configures and open a new connection to the db and returns a [Storage] or an error
func New(cfg *db.Config) (s *Storage, err error) {

	database, err := sql.Open(cfg.DSN().DriverName.String(), cfg.DSN().String())
	if err != nil {
		return nil, err
	}

	err = database.Ping()
	if err != nil {
		database, err = newDB(*cfg) //
		if err != nil {
			return nil, err
		}
	}

	if err = schema.Up(cfg.DSN().DriverName, database); err != nil {
		return
	}

	queries := genDBSQL.New(database)

	return &Storage{
		db:      database,
		queries: queries,
	}, nil
}

// Storage represents an internal implementation of [sql.DB]
type Storage struct {
	db      *sql.DB
	queries *genDBSQL.Queries
}

var _ storage.User = (*Storage)(nil)
var _ storage.Balance = (*Storage)(nil)

// RetrieveUserBalance retrieves current user's balance and the amount withdrawn by their userID or an error
func (s *Storage) RetrieveUserBalance(ctx context.Context, userID uuid.UUID) (currentBalance, withdrawn int, err error) {
	return 0, 0, nil
}

// RetrieveUser retrieves a user id by their login
func (s *Storage) RetrieveUser(ctx context.Context, login string) (userID uuid.UUID, err error) {
	userID, err = s.queries.SelectUserIDByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.UUID{}, storageErrors.ErrNotFound
		}
		return uuid.UUID{}, err
	}
	return userID, nil
}

// StoreUser stores the user by their name and their hashed password.
//   - name MUST be unique
func (s *Storage) StoreUser(ctx context.Context, login, hashedPassword string) (err error) {
	_, err = s.RetrieveUser(ctx, login)
	if err != nil {
		if errors.Is(err, storageErrors.ErrNotFound) {
			goto newUser
		}
		return err
	} else {
		return storageErrors.ErrAlreadyExists
	}

newUser:
	newUserID, err := uuid.NewV7()
	if err != nil {
		return err
	}
	result, err := s.queries.InsertUser(ctx,
		genDBSQL.InsertUserParams{
			ID:             newUserID,
			Login:          login,
			HashedPassword: hashedPassword,
			CreatedAt:      time.Now().UTC(),
			UpdatedAt:      time.Now().UTC()},
	)
	if err != nil {
		return err
	}
	num, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if num != 1 {
		return fmt.Errorf("%w: expected to affect 1 row, affected %d", storageErrors.ErrNoAffect, num)
	}

	return nil
}

func (s *Storage) StoreOrder(ctx context.Context, userID uuid.UUID, orderNumber, orderStatus string, createdAt time.Time) error {

	newOrderID, err := uuid.NewV7()
	if err != nil {
		return err
	}
	res, err := s.queries.InsertOrder(ctx,
		genDBSQL.InsertOrderParams{
			ID:        newOrderID,
			UserID:    userID,
			Number:    orderNumber,
			Status:    orderStatus,
			CreatedAt: createdAt,
		})
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return storageErrors.ErrAlreadyExists
	}

	return nil
}
func (s *Storage) RetreiveOrder(ctx context.Context, userID uuid.UUID, orderNumber string) error {
	// s.queries.Se
	return nil
}

func (s *Storage) RetreiveUserPassword(ctx context.Context, login string) (hashedPassword string, err error) {
	hashedPassword, err = s.queries.SelectUserHashedPasswordByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storageErrors.ErrNotFound
		}
		return "", err
	}
	return hashedPassword, nil
}

func (s *Storage) RetreiveOrderUser(ctx context.Context, orderNumber string) (userID uuid.UUID, err error) {
	userID, err = s.queries.SelectUserIDByOrderNumber(ctx, orderNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.UUID{}, storageErrors.ErrNotFound
		}

		return uuid.UUID{}, err
	}
	return userID, nil
}

func (s *Storage) RetrieaveUserOrders(ctx context.Context, userID uuid.UUID) (userOrders []genDBSQL.SelectOrdersByUserIDRow, err error) {
	rows, err := s.queries.SelectOrdersByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, nil
	}

	return rows, nil
}

func (s *Storage) Retrieve(ctx context.Context, userID uuid.UUID) (genDBSQL.SelectBalanceByUserIDRow, error) {
	userBalance, err := s.queries.SelectBalanceByUserID(ctx, userID)
	if err != nil {
		return genDBSQL.SelectBalanceByUserIDRow{}, err
	}
	return userBalance, nil
}

func (s *Storage) RetrieveOrderIDsForAccrual(ctx context.Context) ([]uuid.UUID, error) {
	statuses := []string{schema.OrderStatusNew, schema.OrderStatusProcessing}

	orderIDs, err := s.queries.SelectOrdersIDsByStatuses(ctx, statuses)
	if err != nil {
		return nil, err
	}

	return orderIDs, nil
}

// StoreUserWithdrawal stores an accrual user transaction
func (s *Storage) RetrieveOrderForAccrual(ctx context.Context, orderID uuid.UUID) (storage.Order, error) {

	order, err := s.queries.SelectOrder(ctx, orderID)
	if err != nil {
		return storage.Order{}, err
	}

	return order, nil
}

func (s *Storage) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status string) error {
	err := s.queries.UpdateOrderStatus(ctx, genDBSQL.UpdateOrderStatusParams{ID: orderID, Status: status})
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) StoreUserAccrual(ctx context.Context, userID uuid.UUID, orderID uuid.UUID, amount int32) error {
	newTransactionID, err := uuid.NewV7()
	if err != nil {
		return nil
	}

	err = s.queries.InsertBalanceTransaction(ctx, genDBSQL.InsertBalanceTransactionParams{
		ID:      newTransactionID,
		Kind:    schema.TransactionKindAccrual,
		UserID:  userID,
		OrderID: orderID,
		Amount:  amount,
	})
	if err != nil {
		return err
	}

	return nil
}

// StoreUserWithdrawal stores a withdrawal user transaction
func (s *Storage) StoreUserWithdrawal(ctx context.Context, userID uuid.UUID, orderID uuid.UUID, amount int32) (uuid.UUID, error) {
	newTransactionID, err := uuid.NewV7()
	if err != nil {
		return uuid.UUID{}, err
	}

	err = s.queries.InsertBalanceTransaction(ctx, genDBSQL.InsertBalanceTransactionParams{
		ID:      newTransactionID,
		Kind:    schema.TransactionKindWithdrawal,
		UserID:  userID,
		OrderID: orderID,
		Amount:  amount,
	})
	if err != nil {
		return uuid.UUID{}, err
	}

	return newTransactionID, nil
}

type Tx struct {
	*storage.Tx
	*storage.Storage
}

// BeginTx is the implementation of [storage.Transacter]. It wraps [database/sql.BeginTx]
func (s *Storage) BeginTx(ctx context.Context) (*storage.Tx, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	sTx := &Storage{
		queries: s.queries.WithTx(tx),
	}

	return &storage.Tx{
		Tx:      tx,
		Balance: sTx,
	}, nil

}

func newDB(cfg db.Config) (*sql.DB, error) {
	ctx := context.Background()
	connector, err := newPostgresConnector(ctx, cfg)
	if err != nil {
		return nil, err
	}
	err = createDB(ctx, connector, cfg.DabaseName)
	if err != nil {
		return nil, err
	}

	return sql.Open(string(cfg.DriverName), cfg.DSN().String())
}

// createDB executes CREATE DATABASE with the dbName parameter using the provided [driver.Connetor]
func createDB(ctx context.Context, conn driver.Connector, dbName string) error {
	q := fmt.Sprintf("CREATE DATABASE %s;", dbName)

	db := sql.OpenDB(conn)
	defer db.Close()

	_, err := db.ExecContext(ctx, q)
	if err != nil {
		return err
	}

	return nil
}

func newPostgresConnector(ctx context.Context, cfg db.Config) (driver.Connector, error) {
	if cfg.DriverName.String() != string(db.DriverNamePostgres) &&
		cfg.DabaseName != string(db.DriverNamePostgres) {
		return nil, storageErrors.ErrUnsupportedDataSource
	}

	connecter, err := pgDriver.NewConnector(cfg.DSN().Default)
	if err != nil {
		return nil, err
	}

	return connecter, nil
}
