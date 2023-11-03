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
		return vehicle
	})

	existing_vehicles, err := v.GetVehicles()
	if err == nil {
		for _, vehicle := range new_vehicles {
			for _, existing_vehicle := range existing_vehicles {
				if vehicle.ID == existing_vehicle.ID {
					vehicle.CreatedOn = existing_vehicle.CreatedOn
					vehicle.UpdatedOn = time.Now()
					break
				}
			}
		}

	}

	db := v.db_service.Connect()
	for _, vehicle := range new_vehicles {

		err := db.Save(&vehicle).Error
		if err != nil {
			continue
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
