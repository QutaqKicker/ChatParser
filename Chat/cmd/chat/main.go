package main

import (
	"chat/internal/config"
	"chat/internal/grpc"
	"context"
	"database/sql"
	"fmt"
	"github.com/QutaqKicker/ChatParser/Common/constants"
	"github.com/QutaqKicker/ChatParser/Common/myKafka"
	_ "github.com/lib/pq"
	"github.com/segmentio/kafka-go"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	cfg := config.MustLoad()

	db, err := connectDb(cfg.Db)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	producer := NewAuditProducer()
	defer producer.Close()

	logger := setupLogger(ctx, producer)

	logger.Info("started application", slog.Any("config", cfg))

	port, _ := strconv.Atoi(os.Getenv(constants.ChatPortEnvName))
	application := grpc.New(logger, db, port)

	go application.Run()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop
	application.Stop()
}

func NewAuditProducer() *myKafka.AuditProducer {
	brokerPort := os.Getenv(constants.KafkaBroker1PortEnvName)

	return &myKafka.AuditProducer{
		kafka.Writer{
			Addr:     kafka.TCP(fmt.Sprintf("localhost:%s", brokerPort)),
			Topic:    constants.KafkaAuditCreateLogTopicName,
			Balancer: &kafka.LeastBytes{}},
	}
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

func setupLogger(ctx context.Context, producer *myKafka.AuditProducer) *slog.Logger {
	log := slog.New(&AuditLogHandler{producer: producer})
	return log
}

type AuditLogHandler struct {
	producer *myKafka.AuditProducer
}

func (h *AuditLogHandler) Handle(ctx context.Context, record slog.Record) error {
	fmt.Println(record.Message)

	err := h.producer.Send(ctx, myKafka.CreateLogRequest{
		ServiceName: "ChatService",
		Type:        int32(record.Level),
		Message:     record.Message,
		Created:     time.Now(),
	})

	if err != nil {
		fmt.Println(err)
	}
	return err
}

func (h *AuditLogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return true
}
func (h *AuditLogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *AuditLogHandler) WithGroup(name string) slog.Handler {
	return h
}
