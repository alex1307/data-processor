package kafka

import (
	"context"
	"errors"
	"log"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	Reader    *kafka.Reader
	Processor MessageProcessor
}

func NewKafkaConsumer(topic string, groupID string, brokers []string, processor MessageProcessor) *KafkaConsumer {
	return &KafkaConsumer{
		Processor: processor,
		Reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			GroupID: groupID,
			Topic:   topic,
		}),
	}
}

func (k *KafkaConsumer) Consume(ctx context.Context) {
	messageBuffer := make([][]byte, 0, 100)
	timeout := time.Duration(5) * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			// Handle the context being canceled
			log.Println("Context canceled, stopping the consumer")
			return

		case <-timer.C:
			// Timeout reached, process messages if there are any
			if len(messageBuffer) > 0 {
				k.Processor.ProcessMessages(messageBuffer)
				messageBuffer = messageBuffer[:0]
			}
			log.Println("Timeout reached, no messages received for", timeout.Seconds(), "seconds")
			timer.Reset(timeout)

		default:
			msg, err := k.Reader.ReadMessage(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					// Context has been canceled, stop processing
					log.Println("Context canceled during message read, stopping the consumer")
					return
				}
				log.Printf("error while receiving message: %s", err.Error())
				continue
			}

			headers := msg.Headers
			for _, header := range headers {
				if header.Key == "command" && string(header.Value) == "stop" {
					return
				}
				if header.Key == "command" && string(header.Value) == "process" {
					messageBuffer = append(messageBuffer, msg.Value)
				}
			}
			// messageBuffer = append(messageBuffer, msg.Value)
			if len(messageBuffer) >= 100 {
				k.Processor.ProcessMessages(messageBuffer)
				messageBuffer = messageBuffer[:0]
				timer.Reset(timeout)
			}
		}
	}
}

func (k *KafkaConsumer) Close() {
	k.Reader.Close()
}
