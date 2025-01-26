package main

import (
	"chat/internal/config"
	"chat/internal/grpc"
	"github.com/QutaqKicker/ChatParser/Common/constants"
	"github.com/QutaqKicker/ChatParser/Common/dbHelper"
	"github.com/QutaqKicker/ChatParser/Common/myKafka"
	"github.com/QutaqKicker/ChatParser/Common/myLogs"
	_ "github.com/lib/pq"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	db, err := dbHelper.ConnectDb(cfg.Db)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	auditProducer := myKafka.NewAuditProducer()
	defer auditProducer.Close()

	userMessageCounterProducer := myKafka.NewUserMessageCounterProducer()
	defer userMessageCounterProducer.Close()

	logger := myLogs.SetupLogger(auditProducer, "ChatService")

	logger.Info("started application", slog.Any("config", cfg))

	port, _ := strconv.Atoi(os.Getenv(constants.ChatPortEnvName))
	application := grpc.New(logger, db, port, userMessageCounterProducer)

	go application.Run()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop
	application.Stop()
}
