package service

import (
	"data-processor/internal/connect"
	csv "data-processor/internal/model/csv"
	dbmodel "data-processor/internal/model/db"

	"log"
)

type SearchService struct {
	db_service connect.Connect
}

func NewSearchService(db_service connect.Connect) *SearchService {
	return &SearchService{
		db_service,
	}
}

func (s *SearchService) FindAllIn(slinks []string) []dbmodel.Searches {
	var records []dbmodel.Searches
	db := s.db_service.Connect()
	db.Where("id IN ?", slinks).Find(&records)
	return records
}

func (s *SearchService) Count() (int64, error) {
	db := s.db_service.Connect()
	var count int64
	err := db.Model(&dbmodel.Searches{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *SearchService) SaveAll(data []csv.SearchMetadata) []dbmodel.Searches {

	var new_records []dbmodel.Searches
	for _, r := range data {
		record := dbmodel.Searches{
			Slink:       r.Slink,
			Timestamp:   r.Timestamp,
			TotalNumber: r.TotalNumber,
			MinPrice:    r.MinPrice,
			MaxPrice:    r.MaxPrice,
			SaleType:    r.SaleType,
		}
		new_records = append(new_records, record)
	}
	log.Println("Saving ", len(new_records), " records...")
	db := s.db_service.Connect()
	db.Save(&new_records)
	return new_records
}

func (s *SearchService) Find(slink string) (dbmodel.Searches, bool) {
	var record dbmodel.Searches
	database := s.db_service.Connect()
	result := database.Where("slink = ?", slink).First(&record)
	if result.Error != nil {
		return dbmodel.Searches{}, false
	}
	return record, true
}
