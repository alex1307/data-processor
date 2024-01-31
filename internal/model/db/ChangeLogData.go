package dbmodel

type ChangeLogData struct {
	ID                  uint64 `gorm:"column:id; primary_key:yes; type: integer; auto_increment:true"`
	ADVERT_ID           string `gorm:"column:advert_id; type: character varying(36); not null"`
	Source              string `gorm:"column:source; ; type: character varying(25); not null"`
	PublishedOn         string `gorm:"column:published_on; type: character varying(36);"`
	LastModifiedOn      string `gorm:"column:last_modified_on; type: character varying(36);"`
	LastModifiedMessage string `gorm:"column:message; type: character varying(100);"`
	DaysInSale          uint32 `gorm:"column:days_in_sale; ; default 0"`
	Sold                bool   `gorm:"column:sold;  default: false"`
	Promoted            bool   `gorm:"column:promoted;  default: false"`
	CreatedOn           string `gorm:"column:created_on; null; type: date; default: now()"`
}
