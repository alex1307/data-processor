package app

import (
	connect "data-processor/internal/connect"
	modelcsv "data-processor/internal/model/csv"
	dbmodel "data-processor/internal/model/db"
	csvservice "data-processor/internal/service/csv"
	dbservice "data-processor/internal/service/db"
	"data-processor/utils"
	"fmt"
	"log"
	"time"

	"github.com/ulule/deepcopier"
)

var (
	vehicle_service   *dbservice.VehicleService
	search_service    *dbservice.SearchService
	equipment_service *dbservice.EquipmentService
	audit_service     *dbservice.AuditService
)

func init() {
	config := connect.GetPosgresConfig("resources/config/postgres_config.yml")
	db_service := connect.GetDBService(config)
	search_service = dbservice.NewSearchService(db_service)
	vehicle_service = dbservice.NewVehicleService(db_service)
	equipment_service = dbservice.NewEquipmentService("resources/config/equipment_config.yml", db_service)
	audit_service = dbservice.NewAuditService(db_service.Connect())
}

func ProcessMetaSearches(file_name string) {
	csv_searches := csvservice.NewGenericCSVReaderService[modelcsv.SearchMetadata]()
	err := csv_searches.ReadFromFiles(file_name)
	if err == nil {
		searches := csv_searches.GetData()
		search_service.SaveAll(searches)
	}
}

func GenerateUpdateList(file_name string, vehicle_service *dbservice.VehicleService) {
	vehicles, err := vehicle_service.GetDataForUpdate()

	if err != nil {
		fmt.Println(err)
		return
	}
	ids := dbservice.Map(vehicles, func(v dbmodel.VehicleRecord) modelcsv.Advert {
		return modelcsv.Advert{ID: v.ID}
	})
	csvservice.WriteToCSVFile(file_name, ids)
}

func ProcessRecords(file_name string, vehicle_service *dbservice.VehicleService) {
	record_service := csvservice.NewRecordService()
	vehicles := record_service.GetRecords([]string{file_name})
	if len(vehicles) > 0 {
		vehicle_service.SaveAll(vehicles)
	}
}

func AuditLog(file_name string, audit_service *dbservice.AuditService) {
	record_service := csvservice.NewRecordService()
	vehicles := record_service.GetRecords([]string{file_name})
	records := dbservice.Map(vehicles, func(source modelcsv.Record) dbmodel.VehicleRecord {
		vehicle := dbmodel.VehicleRecord{}
		deepcopier.Copy(source).To(&vehicle)
		vehicle.CreatedOn = utils.ConvertDate(source.CreatedOn)
		vehicle.UpdatedOn = utils.ConvertDate(source.UpdatedOn)
		return vehicle
	})
	if len(records) > 0 {
		audit_service.AuditAll(records)
	}
}

func ProcessEquipments(file_name string, equipment_service *dbservice.EquipmentService) {
	record_service := csvservice.NewRecordService()
	vehicles := record_service.GetRecords([]string{file_name})
	equipments := dbservice.Map(vehicles, func(source modelcsv.Record) int64 {
		return source.Equipment
	})
	if len(equipments) > 0 {
		inserted := equipment_service.SaveAll(&equipments)
		log.Println("inserted equipments: {}", inserted)
	}
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
	records := dbservice.Map(vehicles, func(source modelcsv.Record) dbmodel.VehicleRecord {
		vehicle := dbmodel.VehicleRecord{}
		deepcopier.Copy(source).To(&vehicle)
		vehicle.CreatedOn = utils.ConvertDate(source.CreatedOn)
		vehicle.UpdatedOn = utils.ConvertDate(source.UpdatedOn)
		return vehicle
	})
	existing_vehicles, _ := vehicle_service.GetVehicles()

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
	csvservice.WriteToCSVFile(ids_file_name, ids)

	log.Println("found records: {}", len(vehicles))

	if len(records) > 0 {
		log.Println("auditing records: {}", len(records))
		audit_service.AuditAll(records)
	}

	if len(vehicles) > 0 {
		vehicle_service.SaveAll(vehicles)
	}

}

func ProcessDeletedIds(data_folder string, file_name string) {
	var data_file_name string
	if file_name == "" {
		data_file_name = fmt.Sprintf("%s/%s-%s.csv", data_folder, "vehicle", time.Now().Format("2006-01-02"))
	} else {
		data_file_name = fmt.Sprintf("%s/%s", data_folder, file_name)
	}
	deleted_service := csvservice.NewDeletedRecordsService()
	deleted := deleted_service.GetDeletedRecords([]string{data_file_name})
	audit_service.LogDeleted(deleted)
}
