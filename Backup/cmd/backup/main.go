package main

import (
	config "backups/internal"
	"backups/internal/grpc"
	"github.com/QutaqKicker/ChatParser/Common/constants"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	logger := setupLogger(cfg.Env)

	logger.Info("started application", slog.Any("config", cfg))

	port, _ := strconv.Atoi(os.Getenv(constants.BackupPortEnvName))
	chatServicePort, _ := strconv.Atoi(os.Getenv(constants.ChatPortEnvName))
	application := grpc.New(logger, cfg.ExportDir, port, chatServicePort)

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
