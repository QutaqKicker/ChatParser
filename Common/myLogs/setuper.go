package myLogs

import (
	"context"
	"fmt"
	"github.com/QutaqKicker/ChatParser/Common/myKafka"
	"log/slog"
	"time"
)

func SetupLogger(producer *myKafka.AuditProducer, serviceName string) *slog.Logger {
	log := slog.New(&AuditLogHandler{producer: producer, serviceName: serviceName})
	return log
}

type AuditLogHandler struct {
	producer    *myKafka.AuditProducer
	serviceName string
}

func (h *AuditLogHandler) Handle(ctx context.Context, record slog.Record) error {
	fmt.Println(record.Message)

	err := h.producer.Send(myKafka.CreateLogRequest{
		ServiceName: h.serviceName,
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
