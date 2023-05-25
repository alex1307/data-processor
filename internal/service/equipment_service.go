package service

import (
	ymlmodel "data-processor/internal/model/yml"
	"log"
	"os"
	"sync"

	yaml "gopkg.in/yaml.v3"
)

var (
	eqservice *EquipmentService
	eqonce    sync.Once
	eqinit    = func(filename string) {
		eqservice = NewEquipmentService(filename)
	}
)

func GetEquipmentService(filename string) *EquipmentService {
	eqonce.Do(func() {
		eqinit(filename)
	})
	return eqservice
}

type EquipmentService struct {
	equipment ymlmodel.Config
}

func NewEquipmentService(filename string) *EquipmentService {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	// Define a map to hold the YAML data
	var equipmentData ymlmodel.Config

	// Unmarshal the YAML data into the struct
	err = yaml.Unmarshal(data, &equipmentData)
	if err != nil {
		log.Fatal(err)
	}
	return &EquipmentService{equipmentData}
}

func (e *EquipmentService) GetEquipment() map[int]string {
	return e.equipment.Equipment
}

func (e *EquipmentService) GetColumns() map[int]string {
	return e.equipment.Mapping
}

func (e *EquipmentService) Equipment2Map(id int32) map[string]bool {

	var indices []int
	for i := 0; id > 0; i++ {
		if id&1 == 1 {
			indices = append(indices, i)
		}
		id >>= 1
	}
	result := map[string]bool{}
	for _, index := range indices {
		column := e.equipment.Mapping[index]
		result[column] = true
	}

	return result
}
