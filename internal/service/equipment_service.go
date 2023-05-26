package service

import (
	dbmodel "data-processor/internal/model/db"
	jsonmodel "data-processor/internal/model/json"
	ymlmodel "data-processor/internal/model/yml"
	"encoding/json"
	"log"
	"os"
	"sync"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/ulule/deepcopier"
	yaml "gopkg.in/yaml.v3"
)

var (
	eqservice *EquipmentService
	eqonce    sync.Once
	eqinit    = func(filename string, db_service *DBService) {
		eqservice = NewEquipmentService(filename, db_service)
	}
)

func GetEquipmentService(filename string, db_service *DBService) *EquipmentService {
	eqonce.Do(func() {
		eqinit(filename, db_service)
	})
	return eqservice
}

type EquipmentService struct {
	equipment  ymlmodel.Config
	db_service *DBService
}

func NewEquipmentService(filename string, db_service *DBService) *EquipmentService {
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
	return &EquipmentService{
		equipmentData,
		db_service,
	}
}

func (e *EquipmentService) GetEquipment() map[int]string {
	return e.equipment.Equipment
}

func (e *EquipmentService) GetColumns() map[int]string {
	return e.equipment.Mapping
}

func (e *EquipmentService) equpipment2map(id int32) map[string]bool {

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

func (s *EquipmentService) Save(EquipmentID int32) int32 {
	equipment := s.from_id(EquipmentID)
	err := s.db_service.db.Save(equipment).Error
	if err != nil {
		log.Println(err)
		return 0
	}
	return equipment.ID
}

func (s *EquipmentService) SaveAll(ListOfEquipmentIds *[]int32) int32 {
	equipments := s.from_slice(ListOfEquipmentIds)
	err := s.db_service.db.Save(equipments).Error
	if err != nil {
		log.Println(err)
		return 0
	}
	return int32(len(equipments))
}

func (s *EquipmentService) from_id(EquipmentID int32) *dbmodel.Equipment {
	target := &dbmodel.Equipment{}
	v := s.equpipment2map(EquipmentID)
	jsonData, _ := json.Marshal(v)
	var expected jsonmodel.EquipmentDTO
	json.Unmarshal(jsonData, &expected)
	deepcopier.Copy(expected).To(target)
	target.ID = EquipmentID
	return target
}

func (s *EquipmentService) from_slice(source *[]int32) []dbmodel.Equipment {
	var records []dbmodel.Equipment = make([]dbmodel.Equipment, 0, len(*source))
	ids := mapset.NewSet[int32]()
	for _, id := range *source {
		equipment := &dbmodel.Equipment{}
		equipment = s.from_id(id)
		if !ids.Contains(id) {
			ids.Add(id)
			records = append(records, *equipment)
		}
	}
	return records
}
