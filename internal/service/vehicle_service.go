package service

import (
	modelcsv "data-processor/internal/model/csv"
	dbmodel "data-processor/internal/model/db"

	"github.com/ulule/deepcopier"
)

type VehicleService struct {
	db_service *DBService
}

func NewVehicleService(db_service *DBService) *VehicleService {
	return &VehicleService{
		db_service,
	}
}

func (v *VehicleService) Save(record modelcsv.Record) (string, error) {
	vehicle := &dbmodel.Vehicle{}
	deepcopier.Copy(record).To(vehicle)
	err := v.db_service.db.Save(vehicle).Error
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
	err := v.db_service.db.Create(&vehicles).Error
	if err != nil {
		return err
	}
	return nil
}
