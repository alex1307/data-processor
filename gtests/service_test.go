package gtests

import (
	"data-processor/internal/connect"
	service "data-processor/internal/service/db"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessCSVFiles(t *testing.T) {
	equipment_config := "../resources/config/equipment_config.yml"
	data_service := service.NewDataService(connect.GetInMemoryConfig(), equipment_config)
	list_files := []string{
		"../resources/test/listing.csv"}
	details_files := []string{
		"../resources/test/details.csv"}
	data_service.ProcessCSVFiles(list_files, details_files)
	assert.True(t, false)
}
