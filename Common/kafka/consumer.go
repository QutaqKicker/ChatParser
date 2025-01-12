package kafka

import (
	"context"
	"fmt"
	"github.com/QutaqKicker/ChatParser/common/constants"
	"github.com/segmentio/kafka-go"
	"os"
)

type AuditConsumer struct {
	*kafka.Reader
}

func NewAuditConsumer() *AuditConsumer {
	brokerUrl := os.Getenv(constants.KafkaBroker1UrlEnvName)

	return &AuditConsumer{
		kafka.NewReader(kafka.ReaderConfig{
			Brokers:  []string{brokerUrl},
			GroupID:  "consumer-group-id",
			Topic:    constants.KafkaAuditCreateLogTopicName,
			MaxBytes: 10e6, // 10MB
		})}
}

func (c *AuditConsumer) ListenRequest(ctx context.Context) (CreateLogRequest, error) {
	message, err := c.ReadMessage(ctx)
	if err != nil {
		return CreateLogRequest{}, err
	}

	request := CreateLogRequest{}
	err = kafka.Unmarshal(message.Value, &request)

	if err != nil {
		return request, err
	}

	fmt.Printf("message at offset %d: %s = %s\n", message.Offset, message.Key, message.Value)
	return request, nil
}

type UserMessageCounterConsumer struct {
	*kafka.Reader
}

func NewUserMessageCounterConsumer() *UserMessageCounterConsumer {
	brokerUrl := os.Getenv(constants.KafkaBroker1UrlEnvName)

	return &UserMessageCounterConsumer{
		kafka.NewReader(kafka.ReaderConfig{
			Brokers:  []string{brokerUrl},
			GroupID:  "consumer-group-id",
			Topic:    constants.KafkaUserMessageCounterTopicName,
			MaxBytes: 10e6, // 10MB
		})}
}

func (c *UserMessageCounterConsumer) ListenRequest(ctx context.Context) (UserMessageCountRequest, error) {
	message, err := c.ReadMessage(ctx)
	if err != nil {
		return UserMessageCountRequest{}, err
	}

	request := UserMessageCountRequest{}
	err = kafka.Unmarshal(message.Value, &request)

	if err != nil {
		return request, err
	}

	fmt.Printf("message at offset %d: %s = %s\n", message.Offset, message.Key, message.Value)
	return request, nil
}
