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
	err := db.Create(&vehicles).Error
	if err != nil {
		return err
	}
	return nil
}
