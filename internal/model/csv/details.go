package modelcsv

type Details struct {
	ID        string `csv:"id"`
	Engine    string `csv:"engine"`
	Gearbox   string `csv:"gearbox"`
	Currency  string `csv:"currency"`
	Price     int32  `csv:"price"`
	Power     int    `csv:"power"`
	Phone     string `csv:"phone"`
	ViewCount int    `csv:"view_count"`
	Equipment int64  `csv:"equipment"`
	CreatedOn string `csv:"created_on"`
}

func (s Details) GetID() string {
	return s.ID
}
