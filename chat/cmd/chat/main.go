package main

import (
	"chat/internal/config"
	db2 "chat/internal/db"
	"chat/internal/domain/models"
	"chat/internal/domain/queryFilters"
	"chat/internal/grpc"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()


	updateValues := make([]db2.UpdateValue, 0)
	updateValues = append(updateValues, db2.UpdateValue{Field: models.Message.})

	fmt.Println(db2.BuildUpdate(queryFilters.MessageFilter{
		SubText: "test",
		ChatIds: []int{1, 2, 3},
		UserIds: []string{"test", "test2"},
	}))

	db, err := connectDb(cfg.Db)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	logger := setupLogger(cfg.Env)

	logger.Info("started application", slog.Any("config", cfg))

	application := grpc.New(logger, db, cfg.Grpc.Port)

	go application.Run()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop
	application.Stop()
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
