package main

import (
	"log/slog"
	"os"
	config "router/internal"
)

func main() {
	cfg := config.MustLoad()

	logger := setupLogger(cfg.Env)

	logger.Info("started application", slog.Any("config", cfg))
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case "dev":
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	}
	return log
}
