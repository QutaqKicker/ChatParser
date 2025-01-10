package kafka

import (
	"github.com/QutaqKicker/ChatParser/common/constants"
	"github.com/segmentio/kafka-go"
	"os"
)

func newAuditProducer() *kafka.Writer {
	brokerUrl := os.Getenv(constants.KafkaBroker1UrlEnvName)

	return &kafka.Writer{
		Addr:     kafka.TCP(brokerUrl),
		Topic:    constants.KafkaAuditCreateLogTopicName,
		Balancer: &kafka.LeastBytes{},
	}
}

func newUserMessageCounterProducer() *kafka.Writer {
	brokerUrl := os.Getenv(constants.KafkaBroker1UrlEnvName)

	return &kafka.Writer{
		Addr:     kafka.TCP(brokerUrl),
		Topic:    constants.KafkaUserMessageCounterTopicName,
		Balancer: &kafka.LeastBytes{},
	}
}
