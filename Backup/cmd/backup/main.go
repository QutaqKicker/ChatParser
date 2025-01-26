package main

import (
	config "backups/internal"
	"backups/internal/grpc"
	"fmt"
	"github.com/QutaqKicker/ChatParser/Common/constants"
	"github.com/QutaqKicker/ChatParser/Common/myKafka"
	"github.com/QutaqKicker/ChatParser/Common/myLogs"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	auditProducer := myKafka.NewAuditProducer()
	defer auditProducer.Close()

	logger := myLogs.SetupLogger(auditProducer, "BackupService")

	logger.Info("started application", slog.Any("config", cfg))

	port, _ := strconv.Atoi(os.Getenv(constants.BackupPortEnvName))
	chatServicePort, _ := strconv.Atoi(os.Getenv(constants.ChatPortEnvName))
	application, err := grpc.New(logger, cfg.ExportDir, port, chatServicePort)
	if err != nil {
		fmt.Println(err)
		return
	}

	go application.Run()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop
	application.Stop()
}
