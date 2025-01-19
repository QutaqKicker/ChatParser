package myKafka

import (
	"context"
	"fmt"
	"github.com/QutaqKicker/ChatParser/Common/constants"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"os"
)

const messagesBatch = 100

type AuditProducer struct {
	kafka.Writer
	messages     []kafka.Message
	messagesChan chan kafka.Message
}

func NewAuditProducer() *AuditProducer {
	brokerUrl := os.Getenv(constants.KafkaBroker1UrlEnvName)

	producer := &AuditProducer{
		Writer: kafka.Writer{
			Addr:     kafka.TCP(brokerUrl),
			Topic:    constants.KafkaAuditCreateLogTopicName,
			Balancer: &kafka.LeastBytes{}},
		messages:     make([]kafka.Message, 0),
		messagesChan: make(chan kafka.Message, 100),
	}

	go producer.startSendingProcess()

	return producer
}

func (p *AuditProducer) startSendingProcess() {
	sendMessagesToKafka := func() {
		err := p.WriteMessages(context.Background(), p.messages...)
		if err != nil {
			fmt.Println(err)
		}
		p.messages = p.messages[:0]
	}
mainLoop:
	for {
		select {
		case message, ok := <-p.messagesChan:
			if !ok {
				break mainLoop
			}
			p.messages = append(p.messages, message)
			if len(p.messages) >= messagesBatch {
				sendMessagesToKafka()
			}
		default:
			if len(p.messages) > 0 {
				sendMessagesToKafka()
			}
		}
	}
}

func (p *AuditProducer) Close() {
	close(p.messagesChan)

	err := p.Writer.Close()
	if err != nil {
		fmt.Println(err)
	}
}

func (p *AuditProducer) Send(message CreateLogRequest) error {
	key := []byte(uuid.New().String())
	value, err := kafka.Marshal(message)
	if err != nil {
		return err
	}

	p.messagesChan <- kafka.Message{
		Key:   key,
		Value: value,
	}

	return err
}

type UserMessageCounterProducer struct {
	kafka.Writer
	messages     []kafka.Message
	messagesChan chan kafka.Message
}

func NewUserMessageCounterProducer() *UserMessageCounterProducer {
	brokerUrl := os.Getenv(constants.KafkaBroker1UrlEnvName)

	producer := &UserMessageCounterProducer{
		Writer: kafka.Writer{
			Addr:     kafka.TCP(brokerUrl),
			Topic:    constants.KafkaUserMessageCounterTopicName,
			Balancer: &kafka.LeastBytes{},
		},
		messages:     make([]kafka.Message, 0),
		messagesChan: make(chan kafka.Message, 100),
	}

	go producer.startSendingProcess()

	return producer
}

func (p *UserMessageCounterProducer) startSendingProcess() {
	sendMessagesToKafka := func() {
		err := p.WriteMessages(context.Background(), p.messages...)
		if err != nil {
			fmt.Println(err)
		}
		p.messages = p.messages[:0]
	}
mainLoop:
	for {
		select {
		case message, ok := <-p.messagesChan:
			if !ok {
				break mainLoop
			}
			p.messages = append(p.messages, message)
			if len(p.messages) >= messagesBatch {
				sendMessagesToKafka()
			}
		default:
			if len(p.messages) > 0 {
				sendMessagesToKafka()
			}
		}
	}
}

func (p *UserMessageCounterProducer) Close() {
	close(p.messagesChan)

	err := p.Writer.Close()
	if err != nil {
		fmt.Println(err)
	}
}

func (p *UserMessageCounterProducer) Send(message UserMessageCountRequest) error {
	key := []byte(uuid.New().String())
	value, err := kafka.Marshal(message)
	if err != nil {
		return err
	}

	p.messagesChan <- kafka.Message{
		Key:   key,
		Value: value,
	}

	return err
}
