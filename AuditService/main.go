package main

import (
	"audit/internal/grpc"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := setupLogger("dev")

	logger.Info("started application")

	port :=
	application := grpc.New(logger, cfg.Grpc.Port)

	go application.Run()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop
	application.Stop()
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
