package dbmodel

type Searches struct {
	Slink       string `gorm:"column:slink; primary_key:yes; type: character varying(36);"`
	Timestamp   int32  `gorm:"column:timestamp; not null; type: integer; default: 0;"`
	TotalNumber int32  `gorm:"column:total_number; not null; type: integer; default: 0;"`
	MinPrice    int32  `gorm:"column:min_price; type: integer; default: 0;"`
	MaxPrice    int32  `gorm:"column:max_price; type: integer; default: 0;"`
	SaleType    string `gorm:"column:deleted_on; type: character varying(25); default: '';"`
}
