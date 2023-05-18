package modelcsv

type NotFound struct {
	ID         string `csv:"id"`
	CreatedOn  string `csv:"created_on"`
	RecordType string `csv:"type"`
}

func (s NotFound) GetID() string {
	return s.ID
}
