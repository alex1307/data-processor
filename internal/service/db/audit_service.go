package service

import (
	"context"
	modelcsv "data-processor/internal/model/csv"
	dbmodel "data-processor/internal/model/db"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const numWorkers = 8 // or whatever number is sensible for your system

type AuditService struct {
	db *gorm.DB
}

func NewAuditService(db_service *gorm.DB) *AuditService {
	return &AuditService{
		db_service,
	}
}

func (a *AuditService) AuditAll(records []dbmodel.VehicleRecord) {
	logrus.Info("auditing changes for all records")
	executeUpdate(records, a.db)
}

func (a *AuditService) LogDeleted(deletedIds []modelcsv.Advert) {
	logrus.Info("auditing changes for all records")
	counter := 0
	for _, deletedId := range deletedIds {
		deleted := dbmodel.DeletedOnAuditLog{
			ADVERT_ID: deletedId.GetID(),
			DeletedOn: time.Now(),
		}
		err := a.db.Save(&deleted).Error
		if condition := err != nil; condition {
			logrus.Error(err)
			continue
		}
		counter++
		if counter%1000 == 0 {
			logrus.Info("deleted ", counter, " records")
			time.Sleep(1 * time.Second)
		}
		record := a.db.Where("id = ?", deletedId.GetID()).Find(&dbmodel.VehicleRecord{})
		if record.Error != nil {
			logrus.Error(record.Error)
			continue
		} else {
			if record.RowsAffected == 0 {
				continue
			}
			deleted_record := dbmodel.VehicleRecord{}
			err := record.First(&deleted_record).Error
			if err != nil {
				logrus.Error(err)
				continue
			}
			deleted_record.DeletedOn = time.Now()
			err = a.db.Save(&deleted_record).Error
			if err != nil {
				logrus.Error(err)
				continue
			}
		}

	}

}

func saveAndAuditRecord(new_record dbmodel.VehicleRecord, db *gorm.DB) {
	old_record := dbmodel.VehicleRecord{}
	err := db.First(&old_record, new_record.ID).Error
	if err != nil {
		db.Save(&new_record)
		return
	}
	if old_record.Price != new_record.Price {
		price_audit_log := dbmodel.PriceAuditLog{
			OldValue:  old_record.Price,
			NewValue:  new_record.Price,
			DiffValue: new_record.Price - old_record.Price,
			ADVERT_ID: new_record.ID,
			UpdatedOn: time.Now(),
		}
		err := db.Save(&price_audit_log)
		if err != nil {
			logrus.Error(err)
		}
	}
	old_date := old_record.UpdatedOn.Format("2006-01-02")
	new_date := new_record.UpdatedOn.Format("2006-01-02")
	if new_date != old_date {
		updated_on_audit_log := dbmodel.UpdatedOnAuditLog{
			OldValue:  old_record.UpdatedOn,
			NewValue:  new_record.UpdatedOn,
			ADVERT_ID: new_record.ID,
			UpdatedOn: time.Now(),
		}
		err := db.Save(&updated_on_audit_log)
		if err != nil {
			logrus.Error(err)
		}
	}

	if old_record.ViewCount != new_record.ViewCount {
		view_count_audit_log := dbmodel.ViewCountAuditLog{
			OldValue:  old_record.ViewCount,
			NewValue:  new_record.ViewCount,
			DiffValue: new_record.ViewCount - old_record.ViewCount,
			ADVERT_ID: new_record.ID,
			UpdatedOn: time.Now(),
		}
		err := db.Save(&view_count_audit_log)
		if err != nil {
			logrus.Error(err)
		}
	}

	if old_record.Equipment != new_record.Equipment {
		equipment_audit_log := dbmodel.EquipmentAuditLog{
			OldValue:  old_record.Equipment,
			NewValue:  new_record.Equipment,
			ADVERT_ID: new_record.ID,
			UpdatedOn: time.Now(),
		}
		err := db.Save(&equipment_audit_log)
		if err != nil {
			logrus.Error(err)
		}
	}

	db.Save(&new_record)

}

func worker(ctx context.Context, id int, records <-chan dbmodel.VehicleRecord, wg *sync.WaitGroup, db *gorm.DB) {
	defer wg.Done()
	for {
		select {
		case record, ok := <-records:
			if !ok {
				return
			}
			saveAndAuditRecord(record, db)
		case <-ctx.Done():
			return
		}
	}
}

func executeUpdate(records []dbmodel.VehicleRecord, db *gorm.DB) {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	recordsChannel := make(chan dbmodel.VehicleRecord)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(ctx, i, recordsChannel, &wg, db)
	}
	counter := 0
	for _, record := range records {
		recordsChannel <- record
		counter++
		if counter%1000 == 0 {
			logrus.Info("audited ", counter, " records")
			time.Sleep(2 * time.Second)
		}
	}

	time.AfterFunc(15*time.Second, func() {
		close(recordsChannel)
		cancel()
	})

	wg.Wait()
}
