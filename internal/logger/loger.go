package logger

import (
	"io"
	"log/slog"
)

func New(w io.Writer) *slog.Logger {
	log := slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{Level: slog.LevelInfo}))
	return log
}
