package kafka

import (
	service "data-processor/internal/service/db"

	"github.com/sirupsen/logrus"
)

type DataProcessor struct {
	protoService service.ProtoService
}

func NewDataProcessor(protoService service.ProtoService) *DataProcessor {
	return &DataProcessor{
		protoService: protoService,
	}
}

func (b *DataProcessor) ProcessMessage(message []byte) error {
	id, err := b.protoService.Save(message)
	if err != nil {
		return err
	}
	logrus.WithFields(logrus.Fields{
		"id": id}).Info("Saved record with id")
	return nil
}

func (b *DataProcessor) ProcessMessages(messages [][]byte) error {
	return b.protoService.SaveAll(messages)
}
