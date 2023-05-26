package service

import (
	csv "data-processor/internal/model/csv"
)

type DataService struct {
	record_service    *RecordService
	equipment_service *EquipmentService
	vehicle_service   *VehicleService
}

func NewDataService(db_config string, equipment_config string) *DataService {
	db_service := NewDBService(db_config)
	return &DataService{
		record_service:    NewRecordService(),
		equipment_service: NewEquipmentService(equipment_config, db_service),
		vehicle_service:   NewVehicleService(db_service),
	}
}

func (s *DataService) ProcessCSVFiles(list_file_names []string, details_file_names []string) error {
	records := s.record_service.GetRecords(list_file_names, details_file_names)
	s.vehicle_service.SaveAll(records)
	equipment_ids := Map(records, func(record csv.Record) int32 {
		return int32(record.Equipment)
	})
	s.equipment_service.SaveAll(&equipment_ids)
	return nil
}

func Map[T any, U any](input []T, mapper func(T) U) []U {
	result := make([]U, len(input))
	for i, v := range input {
		result[i] = mapper(v)
	}
	return result
}
