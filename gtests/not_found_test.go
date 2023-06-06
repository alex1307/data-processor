package gtests

import (
	modelcsv "data-processor/internal/model/csv"
	csvservice "data-processor/internal/service/csv"
	service "data-processor/internal/service/db"
	"log"
	"testing"

	"github.com/elliotchance/pie/v2"
	"github.com/stretchr/testify/assert"
)

var (
	not_found_service = service.NewNotFoundService(db_service)
)

func TestLoadNotFoundErrors(t *testing.T) {
	csv_not_found_service := csvservice.NewGenericCSVReaderService[modelcsv.MobileDataError]()
	errors_filename := "../resources/test/errors.csv"
	err := csv_not_found_service.ReadFromFiles(errors_filename)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	assert.Equal(t, len(csv_not_found_service.GetIdentities().ToSlice()), 420)
	values := pie.Values(csv_not_found_service.GetData())
	first_half := values[0:210]
	assert.Equal(t, len(values), 420)
	saved := not_found_service.SaveAll(first_half)
	counter := 0
	for _, v := range saved {
		log.Println("Saved: ", v)
		if v.Retry == 1 {
			counter++
		}
	}
	assert.Equal(t, saved, 210)
	assert.Equal(t, counter, 210)
	saved = not_found_service.SaveAll(values)

	assert.Equal(t, saved, 420)

	for _, v := range saved {
		if v.Retry == 2 {
			counter++
		}
	}
	assert.Equal(t, counter, 210)
}
