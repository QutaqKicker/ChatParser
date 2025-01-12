package main

import (
	"audit/internal/config"
	"audit/internal/domain/services"
	"context"
	"database/sql"
	"fmt"
	"github.com/QutaqKicker/ChatParser/Common/myKafka"
	_ "github.com/lib/pq"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()
	logger := setupLogger(cfg.Env)

	db, err := connectDb(cfg.Db)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logSaver := services.NewLogSaver(logger, db)
	consumer := myKafka.NewAuditConsumer()
	defer func(consumer *myKafka.AuditConsumer) {
		err := consumer.Close()
		if err != nil {
			panic(err)
		}
	}(consumer)

	requestsChan := consumer.ListenRequests(ctx)

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

				logSaver.SaveLog(r.ServiceName, r.Type, r.Message)
			}

		}
		<-requestsChan
	}()

	logger.Info("started application")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop
}

func connectDb(cfg config.DbConfig) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DbName)
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
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
