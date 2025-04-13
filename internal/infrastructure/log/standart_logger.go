package log

import (
	"log"

	"logur.dev/logur"
)

func NewErrorStandartLogger(logger logur.Logger) *log.Logger {
	return logur.NewErrorStandardLogger(logger, "", 0)
}

func SetStandartLogger(logger logur.Logger) {
	log.SetOutput(logur.NewLevelWriter(logger, logur.Info))
}
