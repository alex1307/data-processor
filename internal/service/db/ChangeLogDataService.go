package service

import (
	"data-processor/internal/connect"
	"data-processor/internal/internal/proto_model"
	dbmodel "data-processor/internal/model/db"
	"log"

	"google.golang.org/protobuf/proto"
)

type ChangeLogDataService struct {
	db_service connect.Connect
}

func NewChangeLogDataService(db_service connect.Connect) ChangeLogDataService {
	return ChangeLogDataService{
		db_service,
	}
}

func (v *ChangeLogDataService) fromBinary(message []byte) (dbmodel.ChangeLogData, error) {
	var source *proto_model.VehicleChangeLogInfo = &proto_model.VehicleChangeLogInfo{}
	proto_err := proto.Unmarshal(message, source)
	if proto_err != nil {
		log.Fatalf("error while unmarshaling protobuf message: %s", proto_err.Error())
		return dbmodel.ChangeLogData{}, proto_err
	}
	log.Println("Unmarshaled protobuf message: ", source)
	var data = dbmodel.ChangeLogData{
		ADVERT_ID:           source.Id,
		Source:              source.Source,
		PublishedOn:         source.PublishedOn,
		LastModifiedOn:      source.LastModifiedOn,
		LastModifiedMessage: source.LastModifiedMessage,
		DaysInSale:          source.DaysInSale,
		Sold:                source.Sold,
		Promoted:            source.Promoted,
	}
	return data, nil
}

func (v ChangeLogDataService) Save(binary []byte) (uint64, error) {
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

func (v ChangeLogDataService) SaveAll(records [][]byte) error {
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

func (v *ChangeLogDataService) GetById(id string) (dbmodel.ChangeLogData, error) {
	vehicle := dbmodel.ChangeLogData{}
	db := v.db_service.Connect()
	err := db.First(&vehicle, id).Error
	if err != nil {
		return vehicle, err
	}
	return vehicle, nil
}

func (v *ChangeLogDataService) GetRecords() ([]dbmodel.ChangeLogData, error) {
	vehicles := []dbmodel.ChangeLogData{}
	db := v.db_service.Connect()
	err := db.Find(&vehicles).Error
	if err != nil {
		return vehicles, err
	}
	return vehicles, nil
}

func (v *ChangeLogDataService) Delete(id string) error {
	vehicle := dbmodel.ChangeLogData{}
	db := v.db_service.Connect()
	err := db.Delete(&vehicle, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (v *ChangeLogDataService) Count() (int64, error) {
	db := v.db_service.Connect()
	var count int64
	err := db.Model(&dbmodel.ChangeLogData{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (v *ChangeLogDataService) FindByListOfIds(ids []string) ([]dbmodel.ChangeLogData, error) {
	vehicles := []dbmodel.ChangeLogData{}
	db := v.db_service.Connect()
	err := db.Where("id in ?", ids).Find(&vehicles).Error
	if err != nil {
		return vehicles, err
	}
	return vehicles, nil
}
