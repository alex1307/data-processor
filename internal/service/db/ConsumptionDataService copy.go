package service

import (
	"data-processor/internal/connect"
	"data-processor/internal/internal/proto_model"
	dbmodel "data-processor/internal/model/db"
	"log"

	"google.golang.org/protobuf/proto"
)

type ConsumptionDataService struct {
	db_service connect.Connect
}

func NewConsumptionDataService(db_service connect.Connect) ConsumptionDataService {
	return ConsumptionDataService{
		db_service,
	}
}

func (v *ConsumptionDataService) fromBinary(message []byte) (dbmodel.ConsumptionData, error) {
	var source *proto_model.Consumption = &proto_model.Consumption{}
	proto_err := proto.Unmarshal(message, source)
	if proto_err != nil {
		log.Fatalf("error while unmarshaling protobuf message: %s", proto_err.Error())
		return dbmodel.ConsumptionData{}, proto_err
	}
	log.Println("Unmarshaled protobuf message: ", source)
	var consumptionData = dbmodel.ConsumptionData{
		ADVERT_ID:       source.Id,
		Source:          source.Source,
		Make:            source.Make,
		Model:           source.Model,
		Year:            source.Year,
		CO2Emission:     source.Co2Emission,
		FuelConsumption: source.FuelConsumption,
		KWConsumption:   source.KwConsuption,
	}
	return consumptionData, nil
}

func (v ConsumptionDataService) Save(binary []byte) (uint64, error) {
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

func (v ConsumptionDataService) SaveAll(records [][]byte) error {
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

func (v *ConsumptionDataService) GetById(id string) (dbmodel.ConsumptionData, error) {
	vehicle := dbmodel.ConsumptionData{}
	db := v.db_service.Connect()
	err := db.First(&vehicle, id).Error
	if err != nil {
		return vehicle, err
	}
	return vehicle, nil
}

func (v *ConsumptionDataService) GetRecords() ([]dbmodel.ConsumptionData, error) {
	vehicles := []dbmodel.ConsumptionData{}
	db := v.db_service.Connect()
	err := db.Find(&vehicles).Error
	if err != nil {
		return vehicles, err
	}
	return vehicles, nil
}

func (v *ConsumptionDataService) Delete(id string) error {
	vehicle := dbmodel.ConsumptionData{}
	db := v.db_service.Connect()
	err := db.Delete(&vehicle, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (v *ConsumptionDataService) Count() (int64, error) {
	db := v.db_service.Connect()
	var count int64
	err := db.Model(&dbmodel.ConsumptionData{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (v *ConsumptionDataService) FindByListOfIds(ids []string) ([]dbmodel.ConsumptionData, error) {
	vehicles := []dbmodel.ConsumptionData{}
	db := v.db_service.Connect()
	err := db.Where("id in ?", ids).Find(&vehicles).Error
	if err != nil {
		return vehicles, err
	}
	return vehicles, nil
}
