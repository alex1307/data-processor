package gtests

import (
	modelcsv "data-processor/internal/model/csv"
	csvservice "data-processor/internal/service/csv"
	service "data-processor/internal/service/db"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestProcessCSVFiles(t *testing.T) {
	time.Sleep(1 * time.Second)
	db := db_service.Connect()
	model := myModel()
	db.Model(&model).Delete(&model)

	records_files := []string{
		"../resources/test/records.csv"}

	record_service := csvservice.NewRecordService()
	vehicles := record_service.GetRecords(records_files)
	time.Sleep(1 * time.Second)
	vehicle_service.SaveAll(vehicles)
	count, _ := vehicle_service.Count()
	assert.Equal(t, int64(67), count)
	count, _ = equipment_service.Count()
	if count > 0 {
		equipment_service.DeleteAll()

	}
	time.Sleep(1 * time.Second)
	count, _ = equipment_service.Count()
	assert.Equal(t, int64(0), count)

	equipment_ids := service.Map(vehicles, func(record modelcsv.Record) int64 {
		return record.Equipment
	})
	assert.Equal(t, 67, len(equipment_ids))
	saved := equipment_service.SaveAll(&equipment_ids)
	time.Sleep(1 * time.Second)
	assert.Equal(t, int32(57), saved)
	count, _ = equipment_service.Count()
	assert.Equal(t, int64(57), count)

}
