package service

import (
	csv "data-processor/internal/model/csv"

	"golang.org/x/exp/maps"
)

type RecordService struct {
	list_service    *GenericCSVReaderService[csv.Listing]
	details_service *GenericCSVReaderService[csv.Details]
}

func NewRecordService() *RecordService {
	return &RecordService{
		list_service:    NewGenericCSVReaderService[csv.Listing](),
		details_service: NewGenericCSVReaderService[csv.Details](),
	}
}

func (s *RecordService) GetRecords(list_file_names []string, details_file_names []string) []csv.Record {
	s.details_service.ReadFromFiles(details_file_names...)
	s.list_service.ReadFromFiles(list_file_names...)
	details_ids := s.details_service.GetIdentities()
	listing_data := s.list_service.GetData()
	details_data := s.details_service.GetData()
	intersection := s.list_service.Intersection(details_ids)
	records := make(map[string]csv.Record, len(intersection.ToSlice()))
	for _, id := range intersection.ToSlice() {
		if _, ok := records[id]; ok {
			continue
		}
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
		records[id] = record
	}
	return maps.Values(records)
}
