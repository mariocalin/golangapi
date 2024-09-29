package log

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

var logger *slog.Logger

func init() {
	w := os.Stderr
	logger = slog.New(tint.NewHandler(w, nil))
}

func Info(msg string) {
	logger.Info(msg)
}

func Debug(msg string) {
	logger.Debug(msg)
}

func Warning(msg string) {
	logger.Warn(msg)
}

func Error(msg string, err error) {
	logger.Error(msg, "error", err)
}

func FatalErr(msg string, err error) {
	logger.Error(msg, "error", err)
	os.Exit(1)
}

func Fatal(msg string) {
	logger.Error(msg)
	os.Exit(1)
}
