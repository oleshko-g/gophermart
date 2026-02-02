// Package db is the base package for db servers
package db

import (
	"net/url"
	"strings"

	storageErrors "github.com/oleshko-g/oggophermart/internal/storage/errors"
)

// Config represents a config of an SQL database
type Config struct {
	dataSource
}

// DSN returns a pointer to the [flag.Value] to set the database source name
func (c *Config) DSN() *dataSource { // revive:disable-line:unexported-return provides the interface to the caller
	return &c.dataSource
}

// dataSource represent a valid Data Source
type dataSource struct {
	_DSN       string
	DabaseName string
	DriverName
	Source  string
	Default string
}

// Set parses s and sets [DSN] and [Driver] or returns an error
func (d *dataSource) Set(s string) error {
	switch d.Source {
	case "":
		d.Source = "env var" // assumes that env var is set upon the first call to [Set]
	case "env var":
		d.Source = "command line flag" // assumes that flags are always parsed last
	}

	url, err := url.Parse(s)
	if err != nil {
		return err
	}

	if url.Scheme != string(DriverNamePostgres) && url.Scheme != string(driverNamePostgreSQL) {
		return storageErrors.ErrUnsupportedDataSource
	}

	d._DSN = url.String()
	// there's only "postgres" SQL driver
	d.DriverName = DriverNamePostgres

	databaseName, _ := strings.CutPrefix(url.Path, "/")
	if databaseName == "" {
		return storageErrors.ErrMissingDatabaseName
	}
	d.DabaseName = databaseName

	d.Default = "postgres://localhost:5432/postgres?sslmode=disable"

	return nil
}

func (d *dataSource) String() string {
	return d._DSN
}

// DriverName is a valid database driver name
type DriverName string

func (d DriverName) String() string {
	return string(d)
}

// Supported database drivers
const (
	DriverNamePostgres   DriverName = "postgres"
	driverNamePostgreSQL DriverName = "postgresql"
	PostgresDefaultDSN              = "postgres://localhost:5432/postgres?sslmode=disable"
)
