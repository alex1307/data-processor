package service

import (
	csv "data-processor/internal/model/csv"
)

type RecordService struct {
	record_service *GenericCSVReaderService[csv.Record]
}

func NewRecordService() *RecordService {
	return &RecordService{
		record_service: NewGenericCSVReaderService[csv.Record](),
	}
}

func (s *RecordService) GetRecords(records_file_names []string) []csv.Record {
	err := s.record_service.ReadFromFiles(records_file_names...)
	if err != nil {
		return []csv.Record{}
	}
	return s.record_service.GetData()
}
