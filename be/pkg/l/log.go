// Package l (log) provides logging functions by wrapping zap.
package l

import (
	"go.uber.org/zap"
)

type Logger struct {
	*zap.Logger
}

func New() Logger {
	logger, err := zapCfg.Build()
	if err != nil {
		panic(err)
	}
	return Logger{logger}
}

func (l Logger) Must(msg string, err error) {
	if err != nil {
		l.Panic(msg, Error(err))
	}
}
