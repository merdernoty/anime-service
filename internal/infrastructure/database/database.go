package database

import (
	"context"
	"database/sql/driver"

	"emperror.dev/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jackc/pgx/v5/tracelog"
)


func NewConnector(config Config) (driver.Connector, error) {
	dsn := config.DSN()

	cfg, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if pgxLogger != nil {
		cfg.Tracer = &tracelog.TraceLog{
			Logger:   pgxLogger,
			LogLevel: tracelog.LogLevelDebug,
		}
	}

	db := stdlib.OpenDB(*cfg)
	defer db.Close()

	if err := db.PingContext(context.Background()); err != nil {
		return nil, errors.WithStack(err)
	}

	return stdlib.GetConnector(*cfg), nil
}
