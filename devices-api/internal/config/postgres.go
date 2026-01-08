package config

import "time"

type PostgresConfig struct {
	URL string

	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration

	StatementTimeout time.Duration
	LockTimeout      time.Duration
}
