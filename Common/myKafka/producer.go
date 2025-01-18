package myKafka

import (
	"context"
	"github.com/QutaqKicker/ChatParser/Common/constants"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"os"
)

type AuditProducer struct {
	kafka.Writer
}

func NewAuditProducer() *AuditProducer {
	brokerUrl := os.Getenv(constants.KafkaBroker1UrlEnvName)

	return &AuditProducer{
		kafka.Writer{
			Addr:     kafka.TCP(brokerUrl),
			Topic:    constants.KafkaAuditCreateLogTopicName,
			Balancer: &kafka.LeastBytes{}},
	}
}

func (p *AuditProducer) Send(ctx context.Context, message CreateLogRequest) error {
	key := []byte(uuid.New().String())
	value, err := kafka.Marshal(message)
	if err != nil {
		return err
	}

	err = p.WriteMessages(ctx, kafka.Message{
		Key:   key,
		Value: value,
	})

	return err
}

type UserMessageCounterProducer struct {
	kafka.Writer
}

func NewUserMessageCounterProducer() *UserMessageCounterProducer {
	brokerUrl := os.Getenv(constants.KafkaBroker1UrlEnvName)

	return &UserMessageCounterProducer{
		kafka.Writer{
			Addr:     kafka.TCP(brokerUrl),
			Topic:    constants.KafkaUserMessageCounterTopicName,
			Balancer: &kafka.LeastBytes{},
		}}
}

func (p *UserMessageCounterProducer) Send(ctx context.Context, message UserMessageCountRequest) error {
	key := []byte(uuid.New().String())
	value, err := kafka.Marshal(message)
	if err != nil {
		return err
	}

	err = p.WriteMessages(ctx, kafka.Message{
		Key:   key,
		Value: value,
	})

	return err
}
