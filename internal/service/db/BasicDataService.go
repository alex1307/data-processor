package service

import (
	"data-processor/internal/connect"
	"data-processor/internal/internal/proto_model"
	dbmodel "data-processor/internal/model/db"
	"log"
	"time"

	"google.golang.org/protobuf/proto"
)

type BasicDataService struct {
	db_service connect.Connect
}

func NewBasicDataService(db_service connect.Connect) BasicDataService {
	return BasicDataService{
		db_service,
	}
}

func (v BasicDataService) fromBinary(message []byte) (dbmodel.BasicData, error) {
	var source *proto_model.BaseVehicleInfo = &proto_model.BaseVehicleInfo{}
	proto_err := proto.Unmarshal(message, source)
	if proto_err != nil {
		log.Fatalf("error while unmarshaling protobuf message: %s", proto_err.Error())
		return dbmodel.BasicData{}, proto_err
	}
	log.Println("Unmarshaled protobuf message: ", source)
	var basicData = dbmodel.BasicData{
		ADVERT_ID: source.Id,
		Source:    source.Source,
		Title:     source.Title,
		Make:      source.Make,
		Model:     source.Model,
		Year:      source.Year,
		Currency:  source.Currency,
		Price:     source.Price,
		Millage:   source.Millage,
		Month:     source.Month,
		Engine:    source.Engine,
		Gearbox:   source.Gearbox,
		CC:        source.Cc,
		PowerPS:   source.PowerPs,
		PowerKW:   source.PowerKw,
	}
	return basicData, nil
}

func (v BasicDataService) Save(binary []byte) (uint64, error) {
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

func (v BasicDataService) SaveAll(records [][]byte) error {
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

func (v *BasicDataService) GetVehicle(id string) (dbmodel.BasicData, error) {
	vehicle := dbmodel.BasicData{}
	db := v.db_service.Connect()
	err := db.First(&vehicle, id).Error
	if err != nil {
		return vehicle, err
	}
	return vehicle, nil
}

func (v *BasicDataService) GetVehicles() ([]dbmodel.BasicData, error) {
	vehicles := []dbmodel.BasicData{}
	db := v.db_service.Connect()
	err := db.Find(&vehicles).Error
	if err != nil {
		return vehicles, err
	}
	return vehicles, nil
}

func (v *BasicDataService) Delete(id string) error {
	vehicle := dbmodel.BasicData{}
	db := v.db_service.Connect()
	err := db.Delete(&vehicle, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (v *BasicDataService) Count() (int64, error) {
	db := v.db_service.Connect()
	var count int64
	err := db.Model(&dbmodel.BasicData{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (v *BasicDataService) GetDataForUpdate() ([]dbmodel.BasicData, error) {
	vehicles := []dbmodel.BasicData{}
	db := v.db_service.Connect()
	today := time.Now()
	day_before := time.Date(today.Year(), today.Month(), today.Day()-1, 0, 0, 0, 0, time.UTC)
	err := db.Where("deleted_on = '0001-01-01' AND created_on <= ? AND YEAR >= 2004", day_before).Find(&vehicles).Error
	if err != nil {
		return vehicles, err
	}

	return vehicles, nil
}

func (v *BasicDataService) FindByListOfIds(ids []string) ([]dbmodel.BasicData, error) {
	vehicles := []dbmodel.BasicData{}
	db := v.db_service.Connect()
	err := db.Where("id in ?", ids).Find(&vehicles).Error
	if err != nil {
		return vehicles, err
	}
	return vehicles, nil
}
