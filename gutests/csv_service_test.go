package gtests

import (
	service "data-processor/internal/service/csv"

	"testing"

	"github.com/stretchr/testify/assert"
)

// func TestFindUniqueIDs(t *testing.T) {
// 	listing_service := service.NewGenericCSVReaderService[modelcsv.Listing]()
// 	listing_filename := "../resources/test/listing.csv"
// 	err := listing_service.ReadFromFiles(listing_filename)
// 	if err != nil {
// 		t.Errorf("Error: %v", err)
// 	}
// 	details_service := service.NewGenericCSVReaderService[modelcsv.Details]()
// 	details_filename := "../resources/test/details.csv"
// 	err = details_service.ReadFromFiles(details_filename)
// 	if err != nil {
// 		t.Errorf("Error: %v", err)
// 	}
// 	listing_ids := listing_service.GetIdentities()
// 	details_ids := details_service.GetIdentities()
// 	assert.Equal(t, listing_ids.Cardinality(), 100)
// 	assert.Equal(t, details_ids.Cardinality(), 81)
// 	intersection := listing_service.Intersection(details_ids)
// 	assert.Equal(t, intersection.Cardinality(), 67)
// 	not_found_details := listing_service.NotFound(details_ids)
// 	assert.Equal(t, not_found_details.Cardinality(), 14)
// 	not_found_listing := details_service.NotFound(listing_ids)
// 	assert.Equal(t, not_found_listing.Cardinality(), 33)
// 	//2 * 67 = 100 + 81 - (33 + 14)
// }

func TestProcess(t *testing.T) {
	records_file_name := "../resources/test/vehicle.csv"
	all, err := service.Process(records_file_name)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	assert.Equal(t, len(all), 67)
}
