package database

import (
	"context"
	"emperror.dev/errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jackc/pgx/v5/tracelog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"logur.dev/logur"
	"time"
)

type GormLogAdapter struct {
	logger logur.Logger
}

func (l *GormLogAdapter) Printf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.logger.Info(msg, nil)
}

func NewConnector(config Config) (*gorm.DB, error) {
	dsn := config.DSN()

	pgxConfig, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if pgxLogger != nil {
		pgxConfig.Tracer = &tracelog.TraceLog{
			Logger:   pgxLogger,
			LogLevel: tracelog.LogLevelDebug,
		}
	}

	sqlDB := stdlib.OpenDB(*pgxConfig)

	if err := sqlDB.PingContext(context.Background()); err != nil {
		sqlDB.Close()
		return nil, errors.WithStack(err)
	}
	var gormLogAdapter *GormLogAdapter
	if pgxLogger != nil && pgxLogger.logger != nil {
		gormLogAdapter = &GormLogAdapter{
			logger: logur.WithField(pgxLogger.logger, "component", "gorm"),
		}
	} else {
		gormLogAdapter = &GormLogAdapter{
			logger: logur.NewNoopLogger(),
		}
	}

	gormLogger := gormlogger.New(
		gormLogAdapter,
		gormlogger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  gormlogger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	gormConfig := &gorm.Config{
		Logger: gormLogger,
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), gormConfig)

	if err != nil {
		sqlDB.Close()
		return nil, errors.WithStack(err)
	}

	return db, nil
}

func CloseDB(db *gorm.DB) error {
	if db == nil {
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		return errors.WithStack(err)
	}
	return sqlDB.Close()
}
 
func AutoMigrate(db *gorm.DB, models ...interface{}) error {
	if db == nil {
		return errors.New("database is nil")
	}

	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}
