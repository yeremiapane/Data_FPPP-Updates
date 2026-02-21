package logger

import (
	"fmt"
	"log/slog"
	"os"
)

var log *slog.Logger

func init() {
	log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}

func Info(msg string, args ...any) {
	log.Info(fmt.Sprintf(msg, args...))
}

func Error(msg string, args ...any) {
	log.Error(fmt.Sprintf(msg, args...))
}
