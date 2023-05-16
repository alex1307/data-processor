package service

import (
	csv "data-processor/internal/model/csv"
	errors_api "errors"
	"fmt"
	"os"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/gocarina/gocsv"
)

type CSVReader interface {
	ReadListings(filename string) ([]*csv.Listing, error)
	ReadDetails(filename string) ([]*csv.Details, error)
	Process(listing_file string, details_file string) ([]*csv.Record, error)
}

type CSVReaderService struct {
	details []csv.Details
	listing []csv.Listing
}

func NewCSVReaderService() *CSVReaderService {
	return &CSVReaderService{}
}

func (s *CSVReaderService) GetListings(filename string) ([]*csv.Listing, error) {
	input, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer input.Close()
	var listings []*csv.Listing
	if err := gocsv.UnmarshalFile(input, &listings); err != nil {
		panic(err)
	}
	return listings, nil
}

func (s *CSVReaderService) GetDetails(filenames ...string) error {
	ids := mapset.NewSet[string]()
	var all_details []csv.Details
	var errors []error
	for _, filename := range filenames {
		input, err := os.Open(filename)
		if err != nil {
			fmt.Println("File not found: ", filename)
			errors = append(errors, err)
			continue
		}
		var details []csv.Details
		if err := gocsv.UnmarshalFile(input, &details); err != nil {
			fmt.Println("Failed to read data from file: ", filename)
			errors = append(errors, err)
			continue
		}
		for _, record := range details {
			if !ids.Contains(record.ID) {
				all_details = append(all_details, record)
			}
			ids.Add(record.ID)
		}
		fmt.Println("Successfully read data from file: ", filename)
		defer input.Close()
	}
	s.details = all_details
	err := joinErrors(errors)
	return err
}

func joinErrors(errors []error) error {
	var errStr string
	for _, err := range errors {
		errStr += err.Error() + "; "
	}
	if errStr == "" {
		return nil
	}
	return errors_api.New(errStr[:len(errStr)-2])
}

func (s *CSVReaderService) Process(listing_file string, details_file string) ([]*csv.Record, error) {
	listings, err := s.GetListings(listing_file)
	if err != nil {
		return nil, err
	}
	details, err := s.GetDetails(details_file)
	if err != nil {
		return nil, err
	}

	var records []*csv.Record
	for _, listing := range listings {
		for _, detail := range details {
			if listing.ID == detail.ID {
				record := csv.Record{
					ID:        listing.ID,
					Model:     listing.Model,
					Millage:   listing.Millage,
					Engine:    detail.Engine,
					Gearbox:   detail.Gearbox,
					Power:     detail.Power,
					Year:      listing.Year,
					Currency:  listing.Currency,
					Price:     listing.Price,
					Phone:     listing.Phone,
					ViewCount: listing.ViewCount,
					Equipment: listing.Equipment,
					Make:      listing.Make,
					Promoted:  listing.Promoted,
					Sold:      listing.Sold,
					Source:    listing.Source,
					CreatedOn: listing.CreatedOn,
				}
				records = append(records, &record)
			}
		}
	}
	return records, nil
}

func (s *CSVReaderService) Intersection(listing []*csv.Listing, details []*csv.Details) ([]*string, error) {
	return nil, nil
}
