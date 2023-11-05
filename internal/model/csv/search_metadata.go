package modelcsv

type SearchMetadata struct {
	Slink       string `csv:"slink"`
	Timestamp   int32  `csv:"timestamp"`
	TotalNumber int32  `csv:"total_number"`
	MinPrice    int32  `csv:"min_price"`
	MaxPrice    int32  `csv:"max_price"`
	SaleType    string `csv:"sale_type"`
}

func (s SearchMetadata) GetID() string {
	return s.Slink
}
