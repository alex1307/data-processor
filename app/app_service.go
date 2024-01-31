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
	"os"

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

func GenerateUpdateList(dir_name string, file_name string) {

	_, err := os.Stat(dir_name)
	if err != nil {
		if os.IsNotExist(err) {

			//create directory
			log.Println("Directory Does not exists and will be created. {}", dir_name)
			os.Mkdir(dir_name, 0755)
		}
	}

	forUpdateFileName := fmt.Sprintf("%s/%s", dir_name, file_name)

	vehicles, err := vehicle_service.GetDataForUpdate()

	if err != nil {
		fmt.Println(err)
		return
	}
	ids := dbservice.Map(vehicles, func(v dbmodel.VehicleRecord) modelcsv.Advert {
		return modelcsv.Advert{ID: v.ID}
	})
	csvservice.WriteToCSVFile(forUpdateFileName, ids)
}

func ProcessRecords(file_name string) {
	record_service := csvservice.NewRecordService()
	vehicles := record_service.GetRecords([]string{file_name})

	var slice_of_100_records = make([]modelcsv.Record, 0, 100)
	var records_for_audit = make([]dbmodel.VehicleRecord, 0, len(vehicles))
	counter := 0
	for _, vehicle := range vehicles {
		slice_of_100_records = append(slice_of_100_records, vehicle)
		counter++
		if counter%100 == 0 {
			for_audit := audit_100_records(slice_of_100_records)
			records_for_audit = append(records_for_audit, for_audit...)
			slice_of_100_records = make([]modelcsv.Record, 0, 100)
		}
	}
	audit_service.AuditAll(records_for_audit)
	if len(vehicles) > 0 {
		vehicle_service.SaveAll(vehicles)
	}
}

func audit_100_records(slice_of_100_records []modelcsv.Record) []dbmodel.VehicleRecord {
	ids := dbservice.Map(slice_of_100_records, func(v modelcsv.Record) string {
		return v.ID
	})
	vehicles, err := vehicle_service.FindByListOfIds(ids)

	if err != nil {
		log.Println("Failed reading data from database: {}", err)
		return []dbmodel.VehicleRecord{}
	}

	for _, vehicle := range vehicles {
		for _, record := range slice_of_100_records {
			if vehicle.ID == record.ID {
				deepcopier.Copy(record).To(&vehicle)
				vehicle.UpdatedOn = utils.ConvertDate(record.UpdatedOn)
			}
		}
	}

	return vehicles
}

func auditLog(file_name string) {
	record_service := csvservice.NewRecordService()
	vehicles := record_service.GetRecords([]string{file_name})
	auditRecords(vehicles)
}

func auditRecords(vehicles []modelcsv.Record) {
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

func processEquipments(file_name string) {
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

func ProcessVehicles(data_dir, source_file_name string, meta_search_file_name string) {
	if meta_search_file_name != "" {
		log.Panicln("Start processing meta searches...")
		ProcessMetaSearches(meta_search_file_name)
	}
	records_file_name := utils.FileName(data_dir, source_file_name)
	auditLog(records_file_name)
	processEquipments(records_file_name)
}

func DeletedVehiclesByIds(data_folder string, file_name string) {
	data_file_name := utils.FileName(data_folder, file_name)
	deleted_service := csvservice.NewDeletedRecordsService()
	deleted := deleted_service.GetDeletedRecords([]string{data_file_name})
	audit_service.LogDeleted(deleted)
}
