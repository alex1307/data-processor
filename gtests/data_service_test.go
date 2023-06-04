package gtests

import (
	service "data-processor/internal/service/db"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessCSVFiles(t *testing.T) {
	equipment_config := "/Users/matkat/Software/Go/src/data-processor/resources/config/equipment_config.yml"
	data_service := service.NewDataService(service.GetInMemoryConfig(), equipment_config)
	list_files := []string{
		"/Users/matkat/Software/release/Rust/scraper/resources/data/listing.csv"}
	details_files := []string{
		"/Users/matkat/Software/release/Rust/scraper/resources/data/details_2023-05-15.csv",
		"/Users/matkat/Software/release/Rust/scraper/resources/data/details_2023-05-16",
		"/Users/matkat/Software/release/Rust/scraper/resources/data/details_2023-05-18",
		"/Users/matkat/Software/release/Rust/scraper/resources/data/details_2023-05-22",
		"/Users/matkat/Software/release/Rust/scraper/resources/data/details_2023-05-23"}
	data_service.ProcessCSVFiles(list_files, details_files)
	assert.True(t, false)
}
