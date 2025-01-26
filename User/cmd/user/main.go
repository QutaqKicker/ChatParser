package main

import (
	"context"
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
	"users/internal/config"
	"users/internal/grpc"
	"users/internal/services"
)

func main() {
	cfg := config.MustLoad()

	db, err := dbHelper.ConnectDb(cfg.Db)
	if err != nil {
		panic(err)
	}

	auditProducer := myKafka.NewAuditProducer()
	defer auditProducer.Close()

	userMessageCounterConsumer := myKafka.NewUserMessageCounterConsumer()
	defer userMessageCounterConsumer.Close()

	logger := myLogs.SetupLogger(auditProducer, "UserService")

	logger.Info("started application", slog.Any("config", cfg))

	port, _ := strconv.Atoi(os.Getenv(constants.UserPortEnvName))
	application := grpc.New(logger, port)

	userMessageCounter := services.NewUserMessageCounter(logger, db)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	requestsChan := userMessageCounterConsumer.ListenRequests(ctx)

	go func() {
	mainLoop:
		for {
			select {
			case <-ctx.Done():
				break mainLoop
			case r, ok := <-requestsChan:
				if !ok {
					break mainLoop
				}

				userMessageCounter.UpdateUserMessagesCount(ctx, r.UserName, int(r.MessageCount))
			}

		}
		<-requestsChan
	}()

	go application.Run()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop
	application.Stop()
}
