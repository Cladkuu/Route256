package logger

import (
	"go.uber.org/zap"
	"os"
)

var (
	GlobalLogger *zap.Logger
)

// logger initiation
func init() {
	// TODO переделать
	logger, err := zap.NewDevelopment()
	if err != nil {
		os.Exit(1)
	}
	GlobalLogger = logger
}
