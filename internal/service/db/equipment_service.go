package service

import (
	"data-processor/internal/connect"
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
	"gorm.io/gorm"
)

var (
	eqservice *EquipmentService
	eqonce    sync.Once
	eqinit    = func(filename string, db_service connect.Connect) {
		eqservice = NewEquipmentService(filename, db_service)
	}
)

func GetEquipmentService(filename string, db_service connect.Connect) *EquipmentService {
	eqonce.Do(func() {
		eqinit(filename, db_service)
	})
	return eqservice
}

type EquipmentService struct {
	equipment  ymlmodel.Config
	db_service connect.Connect
}

func NewEquipmentService(filename string, db_service connect.Connect) *EquipmentService {
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

func (e *EquipmentService) equpipment2map(id int64) map[string]bool {

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

func (s *EquipmentService) Save(EquipmentID int64) int64 {
	equipment := s.FromID(EquipmentID)
	db := s.db_service.Connect()
	err := db.Save(equipment).Error
	if err != nil {
		log.Println(err)
		return 0
	}
	return equipment.ID
}

func (s *EquipmentService) SaveAll(ListOfEquipmentIds *[]int64) int32 {
	equipments := s.FromSlice(ListOfEquipmentIds)
	db := s.db_service.Connect()
	err := db.Save(equipments).Error
	if err != nil {
		log.Println(err)
		return 0
	}
	return int32(len(equipments))
}

func (s *EquipmentService) Delete(EquipmentID int32) int32 {
	db := s.db_service.Connect()
	err := db.Delete(&dbmodel.Equipment{}, EquipmentID).Error
	if err != nil {
		log.Println(err)
		return 0
	}
	return EquipmentID
}

func (s *EquipmentService) Count() (int64, error) {
	db := s.db_service.Connect()
	var count int64
	err := db.Model(&dbmodel.Equipment{}).Count(&count).Error
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return count, nil
}

func (s *EquipmentService) DeleteAll() int32 {
	db := s.db_service.Connect()
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&dbmodel.Equipment{})

	return 1
}

func (s *EquipmentService) FindEquipment(EquipmentID int32) (bool, dbmodel.Equipment) {
	db := s.db_service.Connect()
	var equipment dbmodel.Equipment
	err := db.First(&equipment, EquipmentID).Error
	if err != nil {
		log.Println(err)
		return false, dbmodel.Equipment{}
	}
	return true, equipment
}

func (s *EquipmentService) FromID(EquipmentID int64) *dbmodel.Equipment {
	target := &dbmodel.Equipment{}
	v := s.equpipment2map(EquipmentID)
	jsonData, _ := json.Marshal(v)
	var expected jsonmodel.EquipmentDTO
	json.Unmarshal(jsonData, &expected)
	deepcopier.Copy(expected).To(target)
	target.ID = EquipmentID
	return target
}

func (s *EquipmentService) FromSlice(source *[]int64) []dbmodel.Equipment {
	var records []dbmodel.Equipment = make([]dbmodel.Equipment, 0, len(*source))
	ids := mapset.NewSet[int64]()
	for _, id := range *source {
		equipment := &dbmodel.Equipment{}
		equipment = s.FromID(id)
		if !ids.Contains(id) {
			ids.Add(id)
			records = append(records, *equipment)
		}
	}
	return records
}
