package service

import (
	csv "data-processor/internal/model/csv"
	dbmodel "data-processor/internal/model/db"

	"github.com/ulule/deepcopier"
)

type DataService struct {
	list_service      *GenericCSVReaderService[csv.Listing]
	details_service   *GenericCSVReaderService[csv.Details]
	equipment_service *EquipmentService
	db_service        *DBService
	records           map[string]csv.Record
}

func NewDataService(db_config string, equipment_config string) *DataService {
	return &DataService{
		list_service:      NewGenericCSVReaderService[csv.Listing](),
		details_service:   NewGenericCSVReaderService[csv.Details](),
		equipment_service: NewEquipmentService(equipment_config),
		db_service:        NewDBService(db_config),
		records:           make(map[string]csv.Record),
	}
}

func (s *DataService) loadListData(files ...string) error {
	return s.list_service.ReadFromFiles(files...)
}

func (s *DataService) loadDetailsData(files ...string) error {
	return s.details_service.ReadFromFiles(files...)
}

func (s *DataService) process() {
	details_ids := s.details_service.GetIdentities()
	listing_data := s.list_service.GetData()
	details_data := s.details_service.GetData()
	intersection := s.list_service.Intersection(details_ids)
	for _, id := range intersection.ToSlice() {
		if _, ok := s.records[id]; ok {
			continue
		}
		listing_record := listing_data[id]
		details_record := details_data[id]
		record := csv.Record{
			ID:        listing_record.ID,
			Model:     listing_record.Model,
			Millage:   listing_record.Millage,
			Engine:    details_record.Engine,
			Gearbox:   details_record.Gearbox,
			Power:     details_record.Power,
			Year:      listing_record.Year,
			Currency:  listing_record.Currency,
			Price:     listing_record.Price,
			Phone:     details_record.Phone,
			ViewCount: details_record.ViewCount,
			Equipment: details_record.Equipment,
			Make:      listing_record.Make,
			Promoted:  listing_record.Promoted,
			Sold:      listing_record.Sold,
			Source:    listing_record.Source,
			CreatedOn: listing_record.CreatedOn,
		}
		s.records[id] = record
	}
}

func (s *DataService) vehicles() []dbmodel.Vehicle {
	var records []dbmodel.Vehicle = make([]dbmodel.Vehicle, 0, len(s.records))
	for _, record := range s.records {
		vehicle := &dbmodel.Vehicle{}
		deepcopier.Copy(record).To(vehicle)
		records = append(records, *vehicle)
	}
	return records
}

func (s *DataService) Save() error {
	var errors []error
	var vehicles = s.vehicles()
	var equipments = FromVehicles(&vehicles)
	err := s.db_service.db.Create(&vehicles).Error
	if err != nil {
		errors = append(errors, err)
	}
	err = s.db_service.db.Create(&equipments).Error
	if err != nil {
		errors = append(errors, err)
	}
	return joinErrors(errors)
}

func (s *DataService) ProcessCSVFiles(list_file_names []string, details_file_names []string) error {
	s.loadDetailsData(details_file_names...)
	s.loadListData(list_file_names...)
	s.process()
	err := s.Save()
	if err != nil {
		return err
	}
	return nil
}
