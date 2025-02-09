package main

import (
	"chat/internal/domain/models"
	"chat/internal/grpc"
	"database/sql"
	"github.com/QutaqKicker/ChatParser/Common/config"
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
	cfg := config.MustLoad[models.Config]()

	auditProducer := myKafka.NewAuditProducer()
	defer auditProducer.Close()

	userMessageCounterProducer := myKafka.NewUserMessageCounterProducer()
	defer userMessageCounterProducer.Close()

	logger := myLogs.SetupLogger(auditProducer, "ChatService")

	db, err := dbHelper.ConnectDb(cfg.Db)

	if err != nil {
		panic(err)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logger.Error(err.Error())
		}
	}(db)

	logger.Info("started application", slog.Any("config", cfg))

	port, _ := strconv.Atoi(os.Getenv(constants.ChatPortEnvName))
	application := grpc.New(logger, db, port, userMessageCounterProducer)

	go func() {
		err := application.Run()
		if err != nil {
			logger.Error(err.Error())
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop
	application.Stop()
}
