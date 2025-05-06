// Package pg provides a PostgreSQL implementation of the datastore.Store interface
package pg

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// Option is a function that modifies config
type Option func(*config)

// New creates a new PostgreSQL store with the provided options
func New(opts ...Option) (*Store, error) {
	// Start with default options
	cfg := defaultConfig()

	// Apply provided options
	for _, opt := range opts {
		opt(cfg)
	}

	// Create connection string
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.host, cfg.port, cfg.user, cfg.password, cfg.database, cfg.sslMode,
	)

	// Open database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.maxOpenConns)
	db.SetMaxIdleConns(cfg.maxIdleConns)
	db.SetConnMaxLifetime(cfg.connMaxLife)

	// Test the connection
	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Store{db: db}, nil
}

// NewWithDB creates a new PostgreSQL store with the provided database connection
// This is primarily used for testing
func NewWithDB(db *sql.DB) *Store {
	return &Store{db: db}
}

// WithHost sets the host for the PostgreSQL connection
func WithHost(host string) Option {
	return func(c *config) {
		c.host = host
	}
}

// WithPort sets the port for the PostgreSQL connection
func WithPort(port int) Option {
	return func(c *config) {
		c.port = port
	}
}

// WithUser sets the user for the PostgreSQL connection
func WithUser(user string) Option {
	return func(c *config) {
		c.user = user
	}
}

// WithPassword sets the password for the PostgreSQL connection
func WithPassword(password string) Option {
	return func(c *config) {
		c.password = password
	}
}

// WithDatabase sets the database name for the PostgreSQL connection
func WithDatabase(database string) Option {
	return func(c *config) {
		c.database = database
	}
}

// WithSSLMode sets the SSL mode for the PostgreSQL connection
func WithSSLMode(sslMode string) Option {
	return func(c *config) {
		c.sslMode = sslMode
	}
}

// WithMaxOpenConns sets the maximum number of open connections
func WithMaxOpenConns(maxOpenConns int) Option {
	return func(c *config) {
		c.maxOpenConns = maxOpenConns
	}
}

// WithMaxIdleConns sets the maximum number of idle connections
func WithMaxIdleConns(maxIdleConns int) Option {
	return func(c *config) {
		c.maxIdleConns = maxIdleConns
	}
}

// WithConnMaxLife sets the maximum lifetime of a connection
func WithConnMaxLife(connMaxLife time.Duration) Option {
	return func(c *config) {
		c.connMaxLife = connMaxLife
	}
}
