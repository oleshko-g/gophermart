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
	genDBSQL "github.com/oleshko-g/gophermart/internal/gen/storage/db/sql"
	"github.com/oleshko-g/gophermart/internal/storage"
	"github.com/oleshko-g/gophermart/internal/storage/db"
	"github.com/oleshko-g/gophermart/internal/storage/db/sql/schema"
	storageErrors "github.com/oleshko-g/gophermart/internal/storage/errors"
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

	if err = schema.PostgresUp(database); err != nil {
		return nil, err
	}

	queries := genDBSQL.New(database)

	return &Storage{
		dbName:  cfg.DabaseName,
		db:      database,
		queries: queries,
	}, nil
}

// Storage represents an internal implementation of [sql.DB]
type Storage struct {
	dbName  string
	db      *sql.DB
	queries *genDBSQL.Queries
}

var _ storage.User = (*Storage)(nil)
var _ storage.Balance = (*Storage)(nil)
var _ statementExecer = (*Storage)(nil)

type statementExecer interface {
	Exec(ctx context.Context, stmt string) error
}

// Exec executes the sql statement on the underlying [sql.DB] or returns an error
func (s *Storage) Exec(ctx context.Context, stmt string) error {
	_, err := s.db.ExecContext(ctx, stmt)
	return err
}

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

func (s *Storage) StoreOrder(ctx context.Context, userID uuid.UUID, orderNumber, orderStatus string, createdAt time.Time) (orderID uuid.UUID, err error) {

	newOrderID, err := uuid.NewV7()
	if err != nil {
		return uuid.UUID{}, err
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
		return uuid.UUID{}, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return uuid.UUID{}, err
	}

	if rowsAffected == 0 {
		return uuid.UUID{}, storageErrors.ErrAlreadyExists
	}

	return newOrderID, nil
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

// Retrieve retrieves the user's balance and the amount withdrawn by their userID or an error
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
		ID:        newTransactionID,
		Kind:      schema.TransactionKindAccrual,
		UserID:    userID,
		OrderID:   orderID,
		Amount:    amount,
		CreatedAt: time.Now().UTC(),
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
		ID:        newTransactionID,
		Kind:      schema.TransactionKindWithdrawal,
		UserID:    userID,
		OrderID:   orderID,
		Amount:    amount,
		CreatedAt: time.Now().UTC(),
	})
	if err != nil {
		return uuid.UUID{}, err
	}

	return newTransactionID, nil
}

func (s *Storage) RetrieveUserWithdrawals(ctx context.Context, userID uuid.UUID) (withdrawals []genDBSQL.SelectBalanceOrderTransactionAmountByUserIDAndKindRow, err error) {
	withdrawals, err = s.queries.SelectBalanceOrderTransactionAmountByUserIDAndKind(ctx, genDBSQL.SelectBalanceOrderTransactionAmountByUserIDAndKindParams{
		UserID: userID,
		Kind:   schema.TransactionKindWithdrawal,
	})
	if err != nil {
		return nil, err
	}

	if withdrawals == nil {
		return nil, storageErrors.ErrNotFound
	}

	return withdrawals, nil
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
	connector, err := newPostgresConnector(ctx)
	if err != nil {
		return nil, err
	}
	err = createDB(ctx, connector, cfg.DabaseName)
	if err != nil {
		return nil, err
	}
	return sql.Open(string(cfg.DriverName), cfg.DSN().String())
}

func (s *Storage) TearDown(ctx context.Context) error {
	err := s.db.Close()
	if err != nil {
		return err
	}

	err = dropDB(ctx, s.dbName)
	if err != nil {
		return err
	}

	return nil
}

func dropDB(ctx context.Context, dbName string) error {

	connector, err := newPostgresConnector(ctx)
	if err != nil {
		return err
	}

	db := sql.OpenDB(connector)
	defer db.Close()

	q := fmt.Sprintf("DROP DATABASE %s;", dbName)
	_, err = db.ExecContext(ctx, q)
	if err != nil {
		return err
	}

	return nil
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

func newPostgresConnector(ctx context.Context) (driver.Connector, error) {
	connecter, err := pgDriver.NewConnector(db.PostgresDefaultDSN)
	if err != nil {
		return nil, err
	}

	return connecter, nil
}
