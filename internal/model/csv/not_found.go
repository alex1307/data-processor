package modelcsv

type MobileDataError struct {
	ID        string `csv:"id"`
	Error     string `csv:"error"`
	CreatedOn string `csv:"created_on"`
}

func (s MobileDataError) GetID() string {
	return s.ID
}
