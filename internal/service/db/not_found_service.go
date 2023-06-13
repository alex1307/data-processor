package service

import (
	"data-processor/internal/connect"
	csv "data-processor/internal/model/csv"
	db "data-processor/internal/model/db"
	dbmodel "data-processor/internal/model/db"
	"log"
	"time"
)

type NotFoundService struct {
	db_service connect.Connect
}

func NewNotFoundService(db_service connect.Connect) *NotFoundService {
	return &NotFoundService{
		db_service,
	}
}

func (s *NotFoundService) FindAllIn(ids []string) []db.NotFound {
	var records []db.NotFound
	db := s.db_service.Connect()
	db.Where("id IN ?", ids).Find(&records)
	return records
}

func (s *NotFoundService) Count() (int64, error) {
	db := s.db_service.Connect()
	var count int64
	err := db.Model(&dbmodel.NotFound{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *NotFoundService) SaveAll(data []csv.MobileDataError) []db.NotFound {
	var ids = Map(data, func(record csv.MobileDataError) string {
		return record.ID
	})
	var records = s.FindAllIn(ids)

	var new_records []db.NotFound
	for _, r := range data {
		record := db.NotFound{
			ID:        r.ID,
			Retry:     1,
			CreatedOn: time.Now(),
		}
		for _, found := range records {
			if found.ID == record.ID {
				record.Retry = found.Retry + 1
			}
		}
		new_records = append(new_records, record)

	}
	log.Println("Saving ", len(new_records), " records...")
	db := s.db_service.Connect()
	db.Save(&new_records)
	return new_records
}

func (s *NotFoundService) Find(id string) (db.NotFound, bool) {
	var record db.NotFound
	database := s.db_service.Connect()
	result := database.Where("id = ?", id).First(&record)
	if result.Error != nil {
		return db.NotFound{}, false
	}
	return record, true
}

func (s *NotFoundService) SaveOrUpdate(record csv.MobileDataError) (db.NotFound, bool) {
	dbrecord, found := s.Find(record.ID)
	if found {
		dbrecord.Retry += 1
		dbrecord.UpdatedOn = time.Now()
	} else {
		dbrecord = db.NotFound{
			ID:        record.ID,
			Retry:     1,
			CreatedOn: time.Now(),
		}
	}
	database := s.db_service.Connect()
	result := database.Save(&dbrecord)
	if result.Error != nil {
		return db.NotFound{}, false
	}
	return dbrecord, true
}
