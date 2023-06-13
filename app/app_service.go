package app

import (
	connect "data-processor/internal/connect"
	modelcsv "data-processor/internal/model/csv"
	csvservice "data-processor/internal/service/csv"
	dbservice "data-processor/internal/service/db"
	"data-processor/utils"
	"fmt"

	"github.com/elliotchance/pie/v2"
)

var (
	vehicle_service    *dbservice.VehicleService
	error_service      *dbservice.NotFoundService
	equipment_service  *dbservice.EquipmentService
	csv_errors_service *csvservice.GenericCSVReaderService[modelcsv.MobileDataError]
)

func init() {
	config := connect.GetPosgresConfig("resources/config/postgres_config.yml")
	csv_errors_service = csvservice.NewGenericCSVReaderService[modelcsv.MobileDataError]()
	db_service := connect.GetDBService(config)
	vehicle_service = dbservice.NewVehicleService(db_service)
	error_service = dbservice.NewNotFoundService(db_service)
	equipment_service = dbservice.NewEquipmentService("resources/config/equipment_config.yml", db_service)
}

func ProcessCSVFiles(data_folder string) {
	listing_file := []string{fmt.Sprintf("%s/%s", data_folder, "listing.csv")}
	details_files := utils.ReadFiles(data_folder, "details", "csv")
	error_files := utils.ReadFiles(data_folder, "errors", "csv")
	csv_errors_service.ReadFromFiles(error_files...)
	record_service := csvservice.NewRecordService()
	vehicles := record_service.GetRecords(listing_file, details_files)
	equipment_ids := dbservice.Map(vehicles, func(record modelcsv.Record) int32 {
		return int32(record.Equipment)
	})
	equipment_service.SaveAll(&equipment_ids)
	vehicle_service.SaveAll(vehicles)
	data := csv_errors_service.GetData()
	values := pie.Values(data)
	error_service.SaveAll(values)
}
