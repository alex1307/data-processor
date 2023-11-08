package service

import (
	csv "data-processor/internal/model/csv"
)

type DeletedRecordsService struct {
	deleted_records_service *GenericCSVReaderService[csv.Advert]
}

func NewDeletedRecordsService() *DeletedRecordsService {
	return &DeletedRecordsService{
		deleted_records_service: NewGenericCSVReaderService[csv.Advert](),
	}
}

func (s *DeletedRecordsService) GetDeletedRecords(records_file_names []string) []csv.Advert {
	err := s.deleted_records_service.ReadFromFiles(records_file_names...)
	if err != nil {
		return []csv.Advert{}
	}
	return s.deleted_records_service.GetData()
}
