package service

import (
	csv "data-processor/internal/model/csv"
	utils "data-processor/utils"
	errors_api "errors"
	"fmt"
	"os"
	"sort"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/gocarina/gocsv"
)

type GenericCSVReaderService[T csv.Identity] struct {
	data map[string]T
	keys mapset.Set[string]
}

func NewGenericCSVReaderService[T csv.Identity]() *GenericCSVReaderService[T] {
	return &GenericCSVReaderService[T]{}
}

func (s *GenericCSVReaderService[T]) GetData() map[string]T {
	return s.data
}

func (s *GenericCSVReaderService[T]) GetIdentities() mapset.Set[string] {
	return s.keys
}

func (s *GenericCSVReaderService[T]) Intersection(other mapset.Set[string]) mapset.Set[string] {
	intersection := mapset.NewSet[string]()
	for _, record := range other.ToSlice() {
		if s.keys.Contains(record) {
			intersection.Add(record)
		}
	}
	return intersection
}

func (s *GenericCSVReaderService[T]) NotFound(other mapset.Set[string]) mapset.Set[string] {

	not_found := mapset.NewSet[string]()
	for _, record := range other.ToSlice() {
		if !s.keys.Contains(record) {
			not_found.Add(record)
		}
	}
	return not_found
}

func (s *GenericCSVReaderService[T]) LoadFromSlice(source []T) {
	data := make(map[string]T)
	keys := mapset.NewSet[string]()
	for _, record := range source {
		data[record.GetID()] = record
		keys.Add(record.GetID())
	}
	s.data = data
	s.keys = keys
}

func (s *GenericCSVReaderService[T]) ReadFromFiles(filenames ...string) error {
	ids := mapset.NewSet[string]()
	all_details := make(map[string]T)
	var errors []error
	sort.Sort(utils.DescendingSort(filenames))
	for _, filename := range filenames {
		input, err := os.Open(filename)
		if err != nil {
			fmt.Println("File not found: ", filename)
			errors = append(errors, err)
			continue
		}
		var details []T
		if err := gocsv.UnmarshalFile(input, &details); err != nil {
			fmt.Println("Failed to read data from file: ", filename)
			errors = append(errors, err)
			continue
		}
		for _, record := range details {
			if !ids.Contains(record.GetID()) {
				all_details[record.GetID()] = record
			}
			ids.Add(record.GetID())
		}
		fmt.Println("Successfully read data from file: ", filename)
		defer input.Close()
	}
	s.data = all_details
	s.keys = ids
	err := joinErrors(errors)
	return err
}

func DescendingSort(filenames []string) {
	panic("unimplemented")
}

func (s *GenericCSVReaderService[T]) WriteToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	var details []T
	for _, record := range s.data {
		details = append(details, record)
	}
	return gocsv.MarshalFile(&details, file)
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

func Process(listing_file string, details_file ...string) ([]csv.Record, error) {
	listing := NewGenericCSVReaderService[csv.Listing]()
	details := NewGenericCSVReaderService[csv.Details]()
	err := listing.ReadFromFiles(listing_file)
	if err != nil {
		return nil, err
	}
	err = details.ReadFromFiles(details_file...)
	if err != nil {
		return nil, err
	}

	details_ids := details.GetIdentities()
	listing_data := listing.GetData()
	details_data := details.GetData()
	intersection := listing.Intersection(details_ids)
	var records []csv.Record
	for _, id := range intersection.ToSlice() {
		listing_record := listing_data[id]
		details_record := details_data[id]
		record := csv.Record{
			ID:        listing_record.ID,
			Model:     listing_record.Model,
			Millage:   listing_record.Millage,
			Engine:    details_record.Engine,
			Gearbox:   details_record.Gearbox,
			Power:     details_record.Power,
			Year:      listing_record.Year,
			Currency:  listing_record.Currency,
			Price:     listing_record.Price,
			Phone:     details_record.Phone,
			ViewCount: details_record.ViewCount,
			Equipment: details_record.Equipment,
			Make:      listing_record.Make,
			Promoted:  listing_record.Promoted,
			Sold:      listing_record.Sold,
			Source:    listing_record.Source,
			CreatedOn: listing_record.CreatedOn,
		}
		records = append(records, record)
	}
	return records, nil
}
