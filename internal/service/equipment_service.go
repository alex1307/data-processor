package service

import (
	ymlmodel "data-processor/internal/model/yml"
	"log"
	"os"

	yaml "gopkg.in/yaml.v3"
)

type EquipmentService struct {
	equipment ymlmodel.Equipment
}

func NewEquipmentService(filename string) *EquipmentService {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	// Define a map to hold the YAML data
	var equipmentData ymlmodel.Equipment

	// Unmarshal the YAML data into the struct
	err = yaml.Unmarshal(data, &equipmentData)
	if err != nil {
		log.Fatal(err)
	}
	return &EquipmentService{equipmentData}
}

func (e *EquipmentService) GetEquipment() map[int32]string {
	return e.equipment.Equipment
}

func (e *EquipmentService) GetColumns() map[int32]string {
	return e.equipment.Mapping
}
