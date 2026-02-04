// Package schema is the package to run SQL migrations
package schema

import (
	"database/sql"
	"embed"

	"github.com/oleshko-g/oggophermart/internal/storage/db"
	"github.com/pressly/goose/v3"
)

//go:embed psql/*.sql
var psqlMigrations embed.FS

// PostgresUp runs the schema/psql migration on the database
func PostgresUp(database *sql.DB) error {
	if err := goose.SetDialect(db.DriverNamePostgres.String()); err != nil {
		return err
	}

	var dir string
	goose.SetBaseFS(psqlMigrations)
	dir = "psql"

	if err := goose.Up(database, dir); err != nil {
		return err
	}
	return nil
}

const (
	OrderStatusNew        = "NEW"
	OrderStatusProcessing = "PROCESSING"
	OrderStatusProcessed  = "PROCESSED"
	OrderStatusInvalid    = "INVALID"
)

const (
	TransactionKindAccrual    = "ACCRUAL"
	TransactionKindWithdrawal = "WITHDRAWAL"
)
