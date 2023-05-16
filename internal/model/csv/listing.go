package modelcsv

type Listing struct {
	ID        string `csv:"id"`
	Make      string `csv:"make"`
	Model     string `csv:"model"`
	Currency  string `csv:"currency"`
	Price     int32  `csv:"price"`
	Millage   int    `csv:"millage"`
	Year      int    `csv:"year"`
	Promoted  bool   `csv:"promoted"`
	Sold      bool   `csv:"sold"`
	CreatedOn string `csv:"created_on"`
	Source    string `csv:"dealer"`
}
