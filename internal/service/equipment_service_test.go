package service

import (
	dbmodel "data-processor/internal/model/db"
	jsonmodel "data-processor/internal/model/json"
	"encoding/json"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ulule/deepcopier"
)

var (
	db_service        = NewDBService("../../resources/config/db_config.yml")
	equipment_service = NewEquipmentService("../../resources/config/equipment_config.yml", db_service)
)

func TestEquipmentService_GetEquipment(t *testing.T) {
	equipment := equipment_service.GetEquipment()
	columns := equipment_service.GetColumns()
	assert.Equal(t, len(equipment), len(columns))
	assert.Equal(t, len(equipment), 24)
}

func TestEquipmentService_GetEquipmentName(t *testing.T) {
	v := equipment_service.equpipment2map(24029696)
	for k, v := range v {
		log.Println(k, v)
	}
	jsonData, _ := json.Marshal(v)
	log.Println(string(jsonData))
	var expected jsonmodel.EquipmentDTO
	json.Unmarshal(jsonData, &expected)
	log.Println(expected)
	target := &dbmodel.Equipment{}
	deepcopier.Copy(expected).To(target)
	assert.Equal(t, true, target.CruiseControl)
	assert.Equal(t, true, target.FourWheelDrive)
	assert.Equal(t, true, target.HeatedSeats)
	assert.Equal(t, true, target.FullyServiced)
}
