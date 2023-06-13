package gtests

import (
	modelcsv "data-processor/internal/model/csv"
	csvservice "data-processor/internal/service/csv"

	"testing"

	"github.com/elliotchance/pie/v2"
	"github.com/stretchr/testify/assert"
)

func TestLoadNotFoundErrors(t *testing.T) {
	ResetDB()
	csv_not_found_service := csvservice.NewGenericCSVReaderService[modelcsv.MobileDataError]()
	errors_filename := "../resources/test/errors.csv"
	err := csv_not_found_service.ReadFromFiles(errors_filename)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	assert.Equal(t, len(csv_not_found_service.GetIdentities().ToSlice()), 420)
	values := pie.Values(csv_not_found_service.GetData())
	assert.Equal(t, len(values), 420)
	first_half := values[0:210]
	assert.Equal(t, len(first_half), 210)
	saved := not_found_service.SaveAll(first_half)
	counter := 0
	for _, v := range saved {
		if v.Retry == 1 {
			counter++
		}
	}
	assert.Equal(t, len(saved), 210)
	assert.Equal(t, counter, 210)
	saved = not_found_service.SaveAll(values)

	assert.Equal(t, len(saved), 420)
	counter = 0
	for _, v := range saved {
		if v.Retry == 2 {
			counter++
		}
	}
	assert.Equal(t, counter, 210)
}
