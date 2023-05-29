package service

import (
	"fmt"
	"os/user"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessCSVFiles(t *testing.T) {
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("Failed to get current user:", err)
		return
	}

	// Get the user's home directory
	homeDir := currentUser.HomeDir

	equipment_config := filepath.Join(homeDir, "Software/Go/src/data-processor/resources/config/equipment_config.yml")
	db_config := filepath.Join(homeDir, "Software/Go/src/data-processor/resources/config/db_config.yml")
	data_service := NewDataService(db_config, equipment_config)
	list_files := []string{
		filepath.Join(homeDir, "/Software/release/Rust/scraper/resources/data/listing.csv")}
	details_files := []string{
		filepath.Join(homeDir, "Software/release/Rust/scraper/resources/data/details_2023-05-15.csv"),
		filepath.Join(homeDir, "Software/release/Rust/scraper/resources/data/details_2023-05-16"),
		filepath.Join(homeDir, "Software/release/Rust/scraper/resources/data/details_2023-05-18"),
		filepath.Join(homeDir, "Software/release/Rust/scraper/resources/data/details_2023-05-22"),
		filepath.Join(homeDir, "Software/release/Rust/scraper/resources/data/details_2023-05-29")}
	data_service.ProcessCSVFiles(list_files, details_files)
	assert.True(t, false)
}
