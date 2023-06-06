package service

import (
	"data-processor/internal/connect"
	csv "data-processor/internal/model/csv"
	csvservice "data-processor/internal/service/csv"
)

type DataService struct {
	record_service    *csvservice.RecordService
	equipment_service *EquipmentService
	vehicle_service   *VehicleService
}

func NewDataService(db_config connect.Config, equipment_config string) *DataService {
	db_service := connect.GetDBService(db_config)
	return &DataService{
		record_service:    csvservice.NewRecordService(),
		equipment_service: NewEquipmentService(equipment_config, db_service),
		vehicle_service:   NewVehicleService(db_service),
	}
}

func NewRecordService() {
	panic("unimplemented")
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
