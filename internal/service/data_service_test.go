package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessCSVFiles(t *testing.T) {
	db_config := "/Users/matkat/Software/Go/src/data-processor/resources/config/db_config.yml"
	equipment_config := "/Users/matkat/Software/Go/src/data-processor/resources/config/equipment_config.yml"
	data_service := NewDataService(db_config, equipment_config)
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
