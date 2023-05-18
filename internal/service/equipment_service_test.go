package service

import (
	jsonmodel "data-processor/internal/model/json"
	"encoding/json"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEquipmentService_GetEquipment(t *testing.T) {
	equipment_service := NewEquipmentService("../../resources/config/equipment_mapping.yml")
	equipment := equipment_service.GetEquipment()
	columns := equipment_service.GetColumns()
	assert.Equal(t, len(equipment), len(columns))
	assert.Equal(t, len(equipment), 24)
}

func TestEquipmentService_GetEquipmentName(t *testing.T) {
	equipment_service := NewEquipmentService("../../resources/config/equipment_mapping.yml")
	v := equipment_service.Map2Equipment(24029696)
	for k, v := range v {
		log.Println(k, v)
	}
	jsonData, _ := json.Marshal(v)
	log.Println(string(jsonData))
	var expected jsonmodel.EquipmentDTO
	json.Unmarshal(jsonData, &expected)
	log.Println(expected)
	assert.Equal(t, true, expected.CruiseControl)
	assert.Equal(t, true, expected.FourWheelDrive)
	assert.Equal(t, true, expected.HeatedSeats)
	assert.Equal(t, true, expected.FullyServiced)
}
