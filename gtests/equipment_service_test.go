package gtests

import (
	service "data-processor/internal/service/db"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	db_service        = service.GetDBService(service.GetInMemoryConfig())
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
