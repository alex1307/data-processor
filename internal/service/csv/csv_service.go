package service

import (
	csv "data-processor/internal/model/csv"
	modelcsv "data-processor/internal/model/csv"
	utils "data-processor/utils"
	errors_api "errors"
	"fmt"
	"os"
	"sort"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/gocarina/gocsv"
)

type GenericCSVReaderService[T csv.Identity] struct {
	data []T
	keys mapset.Set[string]
}

func NewGenericCSVReaderService[T csv.Identity]() *GenericCSVReaderService[T] {
	return &GenericCSVReaderService[T]{}
}

func (s *GenericCSVReaderService[T]) GetData() []T {
	return s.data
}

func (s *GenericCSVReaderService[T]) GetIdentities() mapset.Set[string] {
	return s.keys
}

func (s *GenericCSVReaderService[T]) ReadFromFiles(filenames ...string) error {
	var errors []error
	var all_details = []T{}
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
		for _, detail := range details {
			all_details = append(all_details, detail)
		}
		fmt.Println("Successfully read data from file: ", filename)
		defer input.Close()
	}
	s.data = all_details
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

func WriteToCSVFile(filename string, ids []modelcsv.Advert) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return gocsv.MarshalFile(&ids, file)
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

func Process(records_file ...string) ([]csv.Record, error) {
	record_service := NewGenericCSVReaderService[csv.Record]()
	err := record_service.ReadFromFiles(records_file...)
	if err != nil {
		return nil, err
	}
	records := record_service.GetData()
	return records, nil
}
