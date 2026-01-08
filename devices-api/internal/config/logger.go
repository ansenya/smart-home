package config

import (
	"log/slog"
	"os"
)

func initLogger() *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo, // DEBUG / INFO / WARN / ERROR
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger
}
