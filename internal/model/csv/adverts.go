package modelcsv

type Advert struct {
	ID string `csv:"id"`
}

func (s Advert) GetID() string {
	return s.ID
}
