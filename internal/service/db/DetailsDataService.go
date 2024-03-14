package service

import (
	"data-processor/internal/connect"
	"data-processor/internal/internal/proto_model"
	dbmodel "data-processor/internal/model/db"

	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

type DetailsDataService struct {
	db_service connect.Connect
}

func NewDetailsDataService(db_service connect.Connect) DetailsDataService {
	return DetailsDataService{
		db_service,
	}
}

func (v *DetailsDataService) fromBinary(message []byte) (dbmodel.DetailsData, error) {
	var source *proto_model.DetailedVehicleInfo = &proto_model.DetailedVehicleInfo{}
	proto_err := proto.Unmarshal(message, source)
	if proto_err != nil {
		logrus.Error("error while unmarshaling protobuf message: ", proto_err.Error())
		return dbmodel.DetailsData{}, proto_err
	}
	logrus.Info("Unmarshaled protobuf message: ", source)
	var data = dbmodel.DetailsData{
		ADVERT_ID:  source.Id,
		Source:     source.Source,
		Phone:      source.Phone,
		Location:   source.Location,
		ViewCount:  source.ViewCount,
		Equipment:  source.Equipment,
		SellerName: source.SellerName,
		IsDealer:   source.IsDealer,
	}
	return data, nil
}

func (v DetailsDataService) Save(binary []byte) (uint64, error) {
	source, err := v.fromBinary(binary)
	if err != nil {
		return uint64(0), err
	}
	db := v.db_service.Connect()
	logrus.Info("Saving record: ", source)
	result := db.Save(&source)
	if result.Error != nil {
		logrus.Error("error while saving record: ", result.Error.Error())
		return uint64(0), result.Error
	}
	var id = source.ID
	return id, nil
}

func (v DetailsDataService) SaveAll(records [][]byte) error {
	counter := 0
	for _, record := range records {
		source, err := v.Save(record)
		if err != nil {
			logrus.Error("error while saving record: ", source)
			continue
		} else {
			counter++
		}
	}
	logrus.Info("Saved records: ", counter)
	return nil
}

func (v *DetailsDataService) GetById(id string) (dbmodel.DetailsData, error) {
	vehicle := dbmodel.DetailsData{}
	db := v.db_service.Connect()
	err := db.First(&vehicle, id).Error
	if err != nil {
		return vehicle, err
	}
	return vehicle, nil
}

func (v *DetailsDataService) GetRecords() ([]dbmodel.DetailsData, error) {
	vehicles := []dbmodel.DetailsData{}
	db := v.db_service.Connect()
	err := db.Find(&vehicles).Error
	if err != nil {
		return vehicles, err
	}
	return vehicles, nil
}

func (v *DetailsDataService) Delete(id string) error {
	vehicle := dbmodel.DetailsData{}
	db := v.db_service.Connect()
	err := db.Delete(&vehicle, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (v *DetailsDataService) Count() (int64, error) {
	db := v.db_service.Connect()
	var count int64
	err := db.Model(&dbmodel.DetailsData{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (v *DetailsDataService) FindByListOfIds(ids []string) ([]dbmodel.DetailsData, error) {
	vehicles := []dbmodel.DetailsData{}
	db := v.db_service.Connect()
	err := db.Where("id in ?", ids).Find(&vehicles).Error
	if err != nil {
		return vehicles, err
	}
	return vehicles, nil
}
