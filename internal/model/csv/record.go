package modelcsv

type Record struct {
	ID        string `csv:"id"`
	Make      string `csv:"make"`
	Model     string `csv:"model"`
	Millage   int32  `csv:"millage"`
	Engine    string `csv:"engine"`
	Gearbox   string `csv:"gearbox"`
	Power     int32  `csv:"power"`
	Year      int32  `csv:"year"`
	Currency  string `csv:"currency"`
	Price     int32  `csv:"price"`
	Phone     string `csv:"phone"`
	Location  string `csv:"location"`
	ViewCount int32  `csv:"view_count"`
	Equipment int64  `csv:"equipment"`
	Top       bool   `csv:"top"`
	VIP       bool   `csv:"vip"`
	Sold      bool   `csv:"sold"`
	Dealer    bool   `csv:"dealer"`
	CreatedOn string `csv:"created_on"`
	UpdatedOn string `csv:"updated_on"`
}

func (s Record) GetID() string {
	return s.ID
}
