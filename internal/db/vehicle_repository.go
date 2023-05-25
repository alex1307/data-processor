package db

import (
	dbmodel "data-processor/internal/model/db"

	"gorm.io/gorm"
)

type VehicleRepositoryImpl struct {
	db *gorm.DB
}

func NewVehicleRepository(db *gorm.DB) *VehicleRepositoryImpl {
	return &VehicleRepositoryImpl{db: db}
}

func (r *VehicleRepositoryImpl) FindAll() ([]dbmodel.Vehicle, error) {
	var vehicles []dbmodel.Vehicle
	err := r.db.Find(&vehicles).Error
	if err != nil {
		return nil, err
	}
	return vehicles, nil
}

func (r *VehicleRepositoryImpl) Count(filter *dbmodel.Vehicle) (int64, error) {
	var count int64
	err := r.db.Model(filter).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *VehicleRepositoryImpl) Create(vehicle *dbmodel.Vehicle) error {
	err := r.db.Create(vehicle).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *VehicleRepositoryImpl) Update(vehicle *dbmodel.Vehicle) error {
	err := r.db.Save(vehicle).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *VehicleRepositoryImpl) Delete(id int64) error {
	err := r.db.Delete(&dbmodel.Vehicle{}, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *VehicleRepositoryImpl) CreateAll(vehicles *[]dbmodel.Vehicle) (*[]dbmodel.Vehicle, error) {
	err := r.db.Create(vehicles).Error
	if err != nil {
		return nil, err
	}
	return vehicles, nil
}
