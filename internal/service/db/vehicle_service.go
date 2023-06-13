package service

import (
	"data-processor/internal/connect"
	modelcsv "data-processor/internal/model/csv"
	dbmodel "data-processor/internal/model/db"

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
	vehicle := &dbmodel.Vehicle{}
	deepcopier.Copy(record).To(vehicle)
	db := v.db_service.Connect()
	err := db.Save(vehicle).Error
	if err != nil {
		return "", err
	}
	return vehicle.ID, nil
}

func (v *VehicleService) SaveAll(records []modelcsv.Record) error {
	vehicles := Map(records, func(r modelcsv.Record) dbmodel.Vehicle {
		vehicle := dbmodel.Vehicle{}
		deepcopier.Copy(r).To(&vehicle)
		return vehicle
	})

	db := v.db_service.Connect()
	for _, vehicle := range vehicles {

		err := db.Save(&vehicle).Error
		if err != nil {
			continue
		}

	}

	return nil
}

func (v *VehicleService) GetVehicle(id string) (dbmodel.Vehicle, error) {
	vehicle := dbmodel.Vehicle{}
	db := v.db_service.Connect()
	err := db.First(&vehicle, id).Error
	if err != nil {
		return vehicle, err
	}
	return vehicle, nil
}

func (v *VehicleService) GetVehicles() ([]dbmodel.Vehicle, error) {
	vehicles := []dbmodel.Vehicle{}
	db := v.db_service.Connect()
	err := db.Find(&vehicles).Error
	if err != nil {
		return vehicles, err
	}
	return vehicles, nil
}

func (v *VehicleService) Delete(id string) error {
	vehicle := dbmodel.Vehicle{}
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
	err := db.Model(&dbmodel.Vehicle{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
