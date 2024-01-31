package kafka

import (
	service "data-processor/internal/service/db"
	"log"
)

type DataProcessor struct {
	Topic        string
	GroupID      string
	protoService service.ProtoService
}

func NewDataProcessor(topic string, groupID string, protoService service.ProtoService) *DataProcessor {
	return &DataProcessor{
		Topic:        topic,
		GroupID:      groupID,
		protoService: protoService,
	}
}

func (b *DataProcessor) ProcessMessage(message []byte) error {
	id, err := b.protoService.Save(message)
	if err != nil {
		return err
	}
	log.Printf("Saved record with id: %d", id)
	return nil
}

func (b *DataProcessor) ProcessMessages(messages [][]byte) error {
	return b.protoService.SaveAll(messages)
}
