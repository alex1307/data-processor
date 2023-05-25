package service

import (
	dbmodel "data-processor/internal/model/db"
	jsonmodel "data-processor/internal/model/json"
	"encoding/json"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/ulule/deepcopier"
)

var (
	equipment_service = GetEquipmentService("../../resources/config/equipment_config.yml")
)

func FromVehicle(vehicle *dbmodel.Vehicle) *dbmodel.Equipment {
	target := &dbmodel.Equipment{}
	v := equipment_service.Equipment2Map(int32(vehicle.Equipment))
	jsonData, _ := json.Marshal(v)
	var expected jsonmodel.EquipmentDTO
	json.Unmarshal(jsonData, &expected)
	deepcopier.Copy(expected).To(target)
	target.ID = int32(vehicle.Equipment)
	return target
}

func FromVehicles(vehicles *[]dbmodel.Vehicle) []dbmodel.Equipment {
	var records []dbmodel.Equipment = make([]dbmodel.Equipment, 0, len(*vehicles))
	ids := mapset.NewSet[int32]()
	for _, vehicle := range *vehicles {
		equipment := &dbmodel.Equipment{}
		equipment = FromVehicle(&vehicle)
		if !ids.Contains(int32(vehicle.Equipment)) {
			ids.Add(int32(vehicle.Equipment))
			records = append(records, *equipment)
		}
	}
	return records
}
