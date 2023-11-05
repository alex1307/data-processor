package gtests

import (
	"data-processor/internal/connect"
	dbmodel "data-processor/internal/model/db"
	service "data-processor/internal/service/db"
	"testing"
)

var (
	db_service        = connect.GetDBService(connect.GetInMemoryConfig())
	equipment_service = service.NewEquipmentService("../resources/config/equipment_config.yml", db_service)
	vehicle_service   = service.NewVehicleService(db_service)
)

func myModel() []interface{} {
	return []interface{}{
		&dbmodel.Searches{},
		&dbmodel.VehicleRecord{},
		&dbmodel.Equipment{},
	}
}

func ResetDB() {
	db := db_service.Connect()
	model := myModel()
	db.Model(&model).Delete(&model)

}

func TestMain(m *testing.T) {
	ResetDB()
}
