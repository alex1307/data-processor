package gtests

import (
	"data-processor/internal/connect"
	csvservice "data-processor/internal/service/csv"
	service "data-processor/internal/service/db"
	"testing"

	csv "data-processor/internal/model/csv"

	"github.com/stretchr/testify/assert"
)

var (
	db_service        = connect.GetDBService(connect.GetInMemoryConfig())
	equipment_service = service.NewEquipmentService("../resources/config/equipment_config.yml", db_service)
)

func TestEquipmentService_GetEquipment(t *testing.T) {
	equipment := equipment_service.GetEquipment()
	columns := equipment_service.GetColumns()
	assert.Equal(t, len(equipment), len(columns))
	assert.Equal(t, len(equipment), 24)
}

func TestEquipmentService_GetEquipmentName(t *testing.T) {
	v := equipment_service.FromID(24029696)
	assert.Equal(t, true, v.CruiseControl)
	assert.Equal(t, true, v.FourWheelDrive)
	assert.Equal(t, true, v.HeatedSeats)
	assert.Equal(t, true, v.FullyServiced)
	assert.Equal(t, int32(24029696), v.ID)
}

func TestEquipmentService_CRUD(t *testing.T) {
	details_service := csvservice.NewGenericCSVReaderService[csv.Details]()
	details_service.ReadFromFiles([]string{"../resources/test/details.csv"}...)
	details := details_service.GetData()
	assert.Equal(t, 81, len(details))
	var equipment_ids []int32
	for _, detail := range details {
		if detail.Equipment > 0 {
			equipment_ids = append(equipment_ids, int32(detail.Equipment))
		}
	}
	first := equipment_ids[0]
	assert.Equal(t, 79, len(equipment_ids))
	saved := equipment_service.SaveAll(&equipment_ids)
	assert.Equal(t, int32(69), saved)
	found, eq := equipment_service.FindEquipment(first)
	assert.Equal(t, true, found)
	assert.Equal(t, true, eq.NewImport)
	count, _ := equipment_service.Count()
	assert.Equal(t, int64(69), count)
	equipment_service.Delete(first)
	count, _ = equipment_service.Count()
	assert.Equal(t, int64(68), count)
}
