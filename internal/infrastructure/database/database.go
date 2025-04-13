package database

import (
	"context"
	"database/sql"
	"emperror.dev/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jackc/pgx/v5/tracelog"
)

func NewConnector(config Config) (*sql.DB, error) {
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

	if err := db.PingContext(context.Background()); err != nil {
		db.Close()
		return nil, errors.WithStack(err)
	}

	return db, nil
}
