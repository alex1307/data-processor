package gtests

import (
	"context"
	KafkaConsumer "data-processor/internal/kafka"
	"testing"
)

var (
	kafka_consumer = KafkaConsumer.NewKafkaConsumer("test", "test", []string{"127.0.0.1:9094"}, nil)
)

func TestKafkaConsumer(t *testing.T) {
	kafka_consumer.Consume(context.Background())
}
