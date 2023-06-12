package gtests

import (
	csvservice "data-processor/internal/service/csv"
	service "data-processor/internal/service/db"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	vehicle_service = service.NewVehicleService(db_service)
)

func TestVehicleService_CRUD(t *testing.T) {
	record_service := csvservice.NewRecordService()
	list_file_names := []string{"../resources/test/listing.csv"}
	details_file_names := []string{"../resources/test/details.csv"}
	records := record_service.GetRecords(list_file_names, details_file_names)
	assert.Equal(t, 67, len(records))
	vehicle_service.SaveAll(records)
	count, _ := vehicle_service.Count()
	assert.Equal(t, int64(67), count)
	first := records[0]
	vehicle, _ := vehicle_service.GetVehicle(first.ID)
	assert.Equal(t, first.ID, vehicle.ID)
	assert.Equal(t, first.Make, vehicle.Make)
	vehicle_service.Delete(first.ID)
	count, _ = vehicle_service.Count()
	assert.Equal(t, int64(66), count)

}
