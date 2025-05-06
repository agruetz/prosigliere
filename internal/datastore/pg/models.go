// Package pg provides a PostgreSQL implementation of the datastore.Store interface
package pg

import (
	"database/sql"
	"time"
)

// Store implements the datastore.Store interface for PostgreSQL
type Store struct {
	db *sql.DB
}

// config holds the configuration for the PostgreSQL store
type config struct {
	host         string
	port         int
	user         string
	password     string
	database     string
	sslMode      string
	maxOpenConns int
	maxIdleConns int
	connMaxLife  time.Duration
}

// defaultConfig returns the default configuration for the PostgreSQL store
func defaultConfig() *config {
	return &config{
		host:         "localhost",
		port:         5432,
		user:         "postgres",
		password:     "",
		database:     "postgres",
		sslMode:      "disable",
		maxOpenConns: 10,
		maxIdleConns: 5,
		connMaxLife:  time.Minute * 5,
	}
}
