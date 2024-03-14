package service

import (
	"data-processor/internal/connect"
	modelcsv "data-processor/internal/model/csv"
	dbmodel "data-processor/internal/model/db"
	"data-processor/utils"
	"time"

	"github.com/ulule/deepcopier"
)

type VehicleService struct {
	db_service connect.Connect
}

func NewVehicleService(db_service connect.Connect) *VehicleService {
	return &VehicleService{
		db_service,
	}
}

func (v *VehicleService) Save(record modelcsv.Record) (string, error) {
	vehicle := &dbmodel.VehicleRecord{}
	deepcopier.Copy(record).To(vehicle)
	db := v.db_service.Connect()
	err := db.Save(vehicle).Error
	if err != nil {
		return "", err
	}
	return vehicle.ID, nil
}

func (v *VehicleService) SaveAll(records []modelcsv.Record) error {

	new_vehicles := Map(records, func(r modelcsv.Record) dbmodel.VehicleRecord {
		vehicle := dbmodel.VehicleRecord{}
		deepcopier.Copy(r).To(&vehicle)
		vehicle.CreatedOn = utils.ConvertDate(r.CreatedOn)
		vehicle.UpdatedOn = utils.ConvertDate(r.UpdatedOn)
		return vehicle
	})

	db := v.db_service.Connect()
	counter := 0
	for _, vehicle := range new_vehicles {

		err := db.Save(&vehicle).Error
		if err != nil {
			continue
		}
		counter++
		if counter%1000 == 0 {
			time.Sleep(2 * time.Second)
		}
	}
	return nil
}

func (v *VehicleService) GetVehicle(id string) (dbmodel.VehicleRecord, error) {
	vehicle := dbmodel.VehicleRecord{}
	db := v.db_service.Connect()
	err := db.First(&vehicle, id).Error
	if err != nil {
		return vehicle, err
	}
	return vehicle, nil
}

func (v *VehicleService) GetVehicles() ([]dbmodel.VehicleRecord, error) {
	vehicles := []dbmodel.VehicleRecord{}
	db := v.db_service.Connect()
	err := db.Find(&vehicles).Error
	if err != nil {
		return vehicles, err
	}
	return vehicles, nil
}

func (v *VehicleService) Delete(id string) error {
	vehicle := dbmodel.VehicleRecord{}
	db := v.db_service.Connect()
	err := db.Delete(&vehicle, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (v *VehicleService) Count() (int64, error) {
	db := v.db_service.Connect()
	var count int64
	err := db.Model(&dbmodel.VehicleRecord{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (v *VehicleService) GetDataForUpdate() ([]dbmodel.VehicleRecord, error) {
	vehicles := []dbmodel.VehicleRecord{}
	db := v.db_service.Connect()
	today := time.Now()
	day_before := time.Date(today.Year(), today.Month(), today.Day()-1, 0, 0, 0, 0, time.UTC)
	err := db.Where("deleted_on = '0001-01-01' AND created_on <= ? AND YEAR >= 2004", day_before).Find(&vehicles).Error
	if err != nil {
		return vehicles, err
	}

	return vehicles, nil
}

func (v *VehicleService) FindByListOfIds(ids []string) ([]dbmodel.VehicleRecord, error) {
	vehicles := []dbmodel.VehicleRecord{}
	db := v.db_service.Connect()
	err := db.Where("id in ?", ids).Find(&vehicles).Error
	if err != nil {
		return vehicles, err
	}

	return vehicles, nil
}
