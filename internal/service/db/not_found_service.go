package service

import (
	csv "data-processor/internal/model/csv"
	db "data-processor/internal/model/db"
	"time"

	"gorm.io/gorm"
)

type NotFoundService struct {
	db *gorm.DB
}

func NewNotFoundService(db_service Connect) *NotFoundService {
	return &NotFoundService{
		eqservice.db_service.Connect(),
	}
}

func (s *NotFoundService) FindAllIn(ids []string) []db.NotFound {
	var records []db.NotFound
	s.db.Where("id IN ?", ids).Find(&records)
	return records
}

func (s *NotFoundService) Find(id string) (db.NotFound, bool) {
	var record db.NotFound
	result := s.db.Where("id = ?", id).First(&record)
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
	result := s.db.Save(&dbrecord)
	if result.Error != nil {
		return db.NotFound{}, false
	}
	return dbrecord, true
}
