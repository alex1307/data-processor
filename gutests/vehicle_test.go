package gtests

import (
	modelcsv "data-processor/internal/model/csv"
	csvservice "data-processor/internal/service/csv"
	utils "data-processor/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVehicleService_CRUD(t *testing.T) {
	ResetDB()
	record_service := csvservice.NewRecordService()
	list_file_names := []string{"../resources/test/listing.csv"}
	details_file_names := []string{"../resources/test/details.csv"}
	records := record_service.GetRecords(list_file_names, details_file_names)
	assert.Equal(t, 67, len(records))
	filtered := utils.Filter(records, func(r modelcsv.Record) bool {
		return r.ID == "11650623703838993"
	})
	found := filtered[0]
	assert.Equal(t, "11650623703838993", found.ID)
	assert.Equal(t, "Mercedes-Benz", found.Make)
	assert.Equal(t, "2023-05-04", found.CreatedOn)
	vehicle_service.SaveAll(records)
	count, _ := vehicle_service.Count()
	assert.Equal(t, int64(67), count)
	ID := "11650623703838993"
	vehicle, _ := vehicle_service.GetVehicle(ID)
	assert.Equal(t, ID, vehicle.ID)
	assert.Equal(t, "Mercedes-Benz", vehicle.Make)
	assert.Equal(t, "E", vehicle.Model)
	date := utils.ConvertDate("2023-05-04")
	assert.Equal(t, date, vehicle.CreatedOn)
	vehicle_service.Delete(ID)
	count, _ = vehicle_service.Count()
	assert.Equal(t, int64(66), count)

}
