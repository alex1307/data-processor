package app

import (
	connect "data-processor/internal/connect"
	modelcsv "data-processor/internal/model/csv"
	dbmodel "data-processor/internal/model/db"
	csvservice "data-processor/internal/service/csv"
	service "data-processor/internal/service/csv"
	dbservice "data-processor/internal/service/db"
	"fmt"
	"log"
	"time"
)

var (
	vehicle_service   *dbservice.VehicleService
	search_service    *dbservice.SearchService
	equipment_service *dbservice.EquipmentService
)

func init() {
	config := connect.GetPosgresConfig("resources/config/postgres_config.yml")
	db_service := connect.GetDBService(config)
	search_service = dbservice.NewSearchService(db_service)
	vehicle_service = dbservice.NewVehicleService(db_service)
	equipment_service = dbservice.NewEquipmentService("resources/config/equipment_config.yml", db_service)
}

func ProcessCSVFiles(data_folder string, file_name string) {
	metadata_file_name := fmt.Sprintf("%s/%s", data_folder, "meta_data.csv")
	csv_searches := csvservice.NewGenericCSVReaderService[modelcsv.SearchMetadata]()
	err := csv_searches.ReadFromFiles(metadata_file_name)
	if err == nil {
		searches := csv_searches.GetData()
		search_service.SaveAll(searches)
	}
	var data_file_name string
	if file_name == "" {
		data_file_name = fmt.Sprintf("%s/%s-%s.csv", data_folder, "vehicle", time.Now().Format("2006-01-02"))
	} else {
		data_file_name = fmt.Sprintf("%s/%s", data_folder, file_name)
	}

	//records_files := fmt.Sprintf("%s/%s", data_folder, "vehicle-{}.csv")
	log.Println("found vehicles files: {}", data_file_name)
	record_service := csvservice.NewRecordService()
	vehicles := record_service.GetRecords([]string{data_file_name})
	existing_vehicles, err := vehicle_service.GetVehicles()
	if err == nil {
		for _, vehicle := range vehicles {
			for _, existing_vehicle := range existing_vehicles {
				if vehicle.ID == existing_vehicle.ID {
					vehicle.UpdatedOn = time.Now().Format("2006-01-02")
					break
				}
			}
		}
	}

	join := dbservice.Filter(existing_vehicles, func(record dbmodel.VehicleRecord) bool {
		for _, vehicle := range vehicles {
			if vehicle.ID == record.ID {
				return false
			}
		}
		return true
	})
	log.Println("found records: {}", len(join))
	ids := dbservice.Map(join, func(record dbmodel.VehicleRecord) modelcsv.Advert {
		return modelcsv.Advert{ID: record.ID}
	})
	ids_file_name := fmt.Sprintf("%s/%s", data_folder, "for_update.csv")
	service.WriteToCSVFile(ids_file_name, ids)

	log.Println("found records: {}", len(vehicles))
	equipment_ids := dbservice.Map(vehicles, func(record modelcsv.Record) int32 {
		return int32(record.Equipment)
	})
	if len(equipment_ids) > 0 {
		equipment_service.SaveAll(&equipment_ids)
	}

	if len(vehicles) > 0 {
		vehicle_service.SaveAll(vehicles)
	}

}
