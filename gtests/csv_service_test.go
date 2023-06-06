package gtests

import (
	modelcsv "data-processor/internal/model/csv"
	service "data-processor/internal/service/csv"
	"strings"
	"time"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindUniqueIDs(t *testing.T) {
	listing_service := service.NewGenericCSVReaderService[modelcsv.Listing]()
	listing_filename := "../resources/test/listing.csv"
	err := listing_service.ReadFromFiles(listing_filename)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	details_service := service.NewGenericCSVReaderService[modelcsv.Details]()
	details_filename := "../resources/test/details.csv"
	err = details_service.ReadFromFiles(details_filename)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	listing_ids := listing_service.GetIdentities()
	details_ids := details_service.GetIdentities()
	assert.Equal(t, listing_ids.Cardinality(), 100)
	assert.Equal(t, details_ids.Cardinality(), 81)
	intersection := listing_service.Intersection(details_ids)
	assert.Equal(t, intersection.Cardinality(), 67)
	not_found_details := listing_service.NotFound(details_ids)
	assert.Equal(t, not_found_details.Cardinality(), 14)
	not_found_listing := details_service.NotFound(listing_ids)
	assert.Equal(t, not_found_listing.Cardinality(), 33)
	//2 * 67 = 100 + 81 - (33 + 14)
}

func TestProcess(t *testing.T) {
	listing_filename := "resources/test/listing.csv"
	details_filename := "resources/test/details.csv"
	all, err := service.Process(listing_filename, details_filename)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	assert.Equal(t, len(all), 67)
	record_service := service.NewGenericCSVReaderService[modelcsv.Record]()
	record_service.LoadFromSlice(all)
	records_filename := "resources/test/records.csv"
	err = record_service.WriteToFile(records_filename)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestNotFound(t *testing.T) {
	listing_service := service.NewGenericCSVReaderService[modelcsv.Listing]()
	listing_filename := "resources/test/listing.csv"
	err := listing_service.ReadFromFiles(listing_filename)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	details_service := service.NewGenericCSVReaderService[modelcsv.Details]()
	details_filename := "resources/test/details.csv"
	err = details_service.ReadFromFiles(details_filename)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	details_ids := details_service.GetIdentities()
	not_found_details := listing_service.NotFound(details_ids)
	not_found_listing := details_service.NotFound(listing_service.GetIdentities())
	assert.Equal(t, not_found_details.Cardinality(), 14)
	var records []modelcsv.NotFound
	currentTime := time.Now()
	dateString := currentTime.Format("2006-01-02")
	for _, id := range not_found_details.ToSlice() {
		if strings.TrimSpace(id) == "" {
			continue
		}
		record := modelcsv.NotFound{ID: id, CreatedOn: dateString, RecordType: "details"}
		records = append(records, record)
	}
	for _, id := range not_found_listing.ToSlice() {
		if strings.TrimSpace(id) == "" {
			continue
		}
		record := modelcsv.NotFound{ID: id, CreatedOn: dateString, RecordType: "listing"}
		records = append(records, record)
	}
	assert.Equal(t, len(records), 47)
	not_found_service := service.NewGenericCSVReaderService[modelcsv.NotFound]()
	not_found_service.LoadFromSlice(records)
	not_found_filename := "resources/test/not_found.csv"
	err = not_found_service.WriteToFile(not_found_filename)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

}
