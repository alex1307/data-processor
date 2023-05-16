package service

import (
	"fmt"
	"os"
	"testing"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/stretchr/testify/assert"
)

func TestReadCSV(t *testing.T) {
	service := CSVReaderService{}
	listing_filename := "/Users/ayagasha/Software/release/Rust/scraper/resources/data/listing.csv"
	listings, err := service.GetListings(listing_filename)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if len(listings) == 0 {
		t.Errorf("No records found")
	}
	assert.Equal(t, len(listings), 3919)
	details_filename := "/Users/ayagasha/Software/release/Rust/scraper/resources/data/details_2023-05-15.csv"
	details, err := service.GetListings(details_filename)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if len(listings) == 0 {
		t.Errorf("No records found")
	}
	assert.Equal(t, len(details), 2236)
}

func TestFindUniqueIDs(t *testing.T) {
	service := CSVReaderService{}
	wd, err := os.Getwd()
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	fmt.Println("Working directory:", wd)
	listing_filename := "../../resources/test/listing.csv"
	listings, err := service.GetListings(listing_filename)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	details_filename := "../../resources/test/details.csv"
	details, err := service.GetListings(details_filename)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	listing_ids := mapset.NewSet[string]()
	details_ids := mapset.NewSet[string]()
	for _, listing := range listings {
		listing_ids.Add(listing.ID)
	}
	assert.Equal(t, listing_ids.Cardinality(), 3718)
	for _, detail := range details {
		details_ids.Add(detail.ID)
	}

	assert.Equal(t, details_ids.Cardinality(), 2232)

	invalid_ids := listing_ids.Difference(details_ids)
	assert.Equal(t, invalid_ids.Cardinality(), 1487)
	intersect := listing_ids.Intersect(details_ids)
	assert.Equal(t, intersect.Cardinality(), 2231)

}
