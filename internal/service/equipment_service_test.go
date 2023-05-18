package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEquipmentService_GetEquipment(t *testing.T) {
	equipment_service := NewEquipmentService("../../resources/config/equipment_mapping.yml")
	equipment := equipment_service.GetEquipment()
	columns := equipment_service.GetColumns()
	assert.Equal(t, len(equipment), len(columns))
	assert.Equal(t, len(equipment), 24)
}
