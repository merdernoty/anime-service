package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/tracelog"
	"logur.dev/logur"
)

type pgxLogurLogger struct {
	logger logur.Logger
}

func (l *pgxLogurLogger) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]any) {
	fields := map[string]interface{}{}
	for k, v := range data {
		fields[k] = fmt.Sprintf("%v", v)
	}

	log := logur.WithFields(l.logger, fields)

	switch level {
	case tracelog.LogLevelTrace, tracelog.LogLevelDebug:
		log.Debug(msg)
	case tracelog.LogLevelInfo:
		log.Info(msg)
	case tracelog.LogLevelWarn:
		log.Warn(msg)
	case tracelog.LogLevelError:
		log.Error(msg)
	default:
		log.Info(msg)
	}
}
var pgxLogger *pgxLogurLogger

func SetLogger(logger logur.Logger) {
	pgxLogger = &pgxLogurLogger{
		logger: logur.WithField(logger, "component", "postgres"),
	}
}

