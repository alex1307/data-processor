package modelcsv

type Record struct {
	ID        string `csv:"id"`
	Model     string `csv:"model"`
	Millage   int    `csv:"millage"`
	Engine    string `csv:"engine"`
	Gearbox   string `csv:"gearbox"`
	Power     int    `csv:"power"`
	Year      int    `csv:"year"`
	Currency  string `csv:"currency"`
	Price     int32  `csv:"price"`
	Phone     string `csv:"phone"`
	ViewCount int    `csv:"view_count"`
	Equipment int64  `csv:"equipment"`
	Make      string `csv:"make"`
	Promoted  bool   `csv:"promoted"`
	Sold      bool   `csv:"sold"`
	Source    string `csv:"dealer"`
	CreatedOn string `csv:"created_on"`
}
