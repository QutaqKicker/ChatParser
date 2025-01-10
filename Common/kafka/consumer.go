package kafka

import (
	"github.com/QutaqKicker/ChatParser/common/constants"
	"github.com/segmentio/kafka-go"
	"os"
)

func newAuditConsumer() *kafka.Reader {
	brokerUrl := os.Getenv(constants.KafkaBroker1UrlEnvName)

	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{brokerUrl},
		GroupID:  "consumer-group-id",
		Topic:    constants.KafkaAuditCreateLogTopicName,
		MaxBytes: 10e6, // 10MB
	})
}

func newUserMessageCounterConsumer() *kafka.Reader {
	brokerUrl := os.Getenv(constants.KafkaBroker1UrlEnvName)

	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{brokerUrl},
		GroupID:  "consumer-group-id",
		Topic:    constants.KafkaUserMessageCounterTopicName,
		MaxBytes: 10e6, // 10MB
	})
}
