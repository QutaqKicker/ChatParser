package myKafka

import (
	"context"
	"fmt"
	"github.com/QutaqKicker/ChatParser/common/pkg/constants"
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

func (c *AuditConsumer) ListenRequests(ctx context.Context) <-chan CreateLogRequest {
	return listenRequests[*AuditConsumer, CreateLogRequest](c, ctx)
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

func (c *UserMessageCounterConsumer) ListenRequests(ctx context.Context) <-chan UserMessageCountRequest {
	return listenRequests[*UserMessageCounterConsumer, UserMessageCountRequest](c, ctx)
}

type kafkaReader interface {
	ReadMessage(ctx context.Context) (kafka.Message, error)
}

func listenRequests[R kafkaReader, T any](reader R, ctx context.Context) <-chan T {
	outChan := make(chan T, 10)

	go func() {
	mainLoop:
		for {
			select {
			case <-ctx.Done():
				break mainLoop
			default:
				message, err := reader.ReadMessage(ctx)
				if err != nil {
					fmt.Println(err)
					break mainLoop
				}

				request := new(T)
				err = kafka.Unmarshal(message.Value, request)
				if err != nil {
					fmt.Println(err)
					break mainLoop
				}

				fmt.Printf("message at offset %d: %s = %s\n", message.Offset, message.Key, message.Value)

				outChan <- *request
			}
		}

		close(outChan)
	}()
	return outChan
}
