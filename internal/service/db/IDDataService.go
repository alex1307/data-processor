package service

import (
	"data-processor/internal/connect"
	"data-processor/internal/internal/proto_model"
	dbmodel "data-processor/internal/model/db"
	"log"
	"time"

	"google.golang.org/protobuf/proto"
)

type IDDataService struct {
	db_service connect.Connect
}

func NewIDDataService(db_service connect.Connect) IDDataService {
	return IDDataService{
		db_service,
	}
}

func (v *IDDataService) fromBinary(message []byte) (dbmodel.IDData, error) {
	var source *proto_model.ID = &proto_model.ID{}
	proto_err := proto.Unmarshal(message, source)
	if proto_err != nil {
		log.Fatalf("error while unmarshaling protobuf message: %s", proto_err.Error())
		return dbmodel.IDData{}, proto_err
	}
	log.Println("Unmarshaled protobuf message: ", source)
	var IDData = dbmodel.IDData{
		ADVERT_ID: source.Id,
		Source:    source.Source,
	}
	return IDData, nil
}

func (v IDDataService) Save(binary []byte) (uint64, error) {
	source, err := v.fromBinary(binary)
	if err != nil {
		return uint64(0), err
	}
	db := v.db_service.Connect()
	log.Println("Saving record: ", source)
	result := db.Save(&source)
	if result.Error != nil {
		log.Fatalf("error while saving record: %s", result.Error.Error())
		return uint64(0), result.Error
	}
	var id = source.ID
	return id, nil
}

func (v IDDataService) SaveAll(records [][]byte) error {
	counter := 0
	for _, record := range records {
		source, err := v.Save(record)
		if err != nil {
			log.Fatalf("error while saving record: %v", source)
			continue
		} else {
			counter++
		}
	}
	log.Printf("Saved %v records", counter)
	return nil
}

func (v *IDDataService) GetVehicle(id string) (dbmodel.IDData, error) {
	vehicle := dbmodel.IDData{}
	db := v.db_service.Connect()
	err := db.First(&vehicle, id).Error
	if err != nil {
		return vehicle, err
	}
	return vehicle, nil
}

func (v *IDDataService) GetVehicles() ([]dbmodel.IDData, error) {
	vehicles := []dbmodel.IDData{}
	db := v.db_service.Connect()
	err := db.Find(&vehicles).Error
	if err != nil {
		return vehicles, err
	}
	return vehicles, nil
}

func (v *IDDataService) Delete(id string) error {
	vehicle := dbmodel.IDData{}
	db := v.db_service.Connect()
	err := db.Delete(&vehicle, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (v *IDDataService) Count() (int64, error) {
	db := v.db_service.Connect()
	var count int64
	err := db.Model(&dbmodel.IDData{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (v *IDDataService) GetDataForUpdate() ([]dbmodel.IDData, error) {
	vehicles := []dbmodel.IDData{}
	db := v.db_service.Connect()
	today := time.Now()
	day_before := time.Date(today.Year(), today.Month(), today.Day()-1, 0, 0, 0, 0, time.UTC)
	err := db.Where("deleted_on = '0001-01-01' AND created_on <= ? AND YEAR >= 2004", day_before).Find(&vehicles).Error
	if err != nil {
		return vehicles, err
	}

	return vehicles, nil
}

func (v *IDDataService) FindByListOfIds(ids []string) ([]dbmodel.IDData, error) {
	vehicles := []dbmodel.IDData{}
	db := v.db_service.Connect()
	err := db.Where("id in ?", ids).Find(&vehicles).Error
	if err != nil {
		return vehicles, err
	}
	return vehicles, nil
}
