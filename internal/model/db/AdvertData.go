package dbmodel

import "time"

type AdvertData struct {
	ID              int32     `gorm:"column:id; primary_key:yes; type: integer; auto_increment:true"`
	ADVERT_ID       string    `gorm:"column:advert_id; type: character varying(36); not null"`
	Source          string    `gorm:"column:source; type: character varying(25); not null"`
	PublishedOn     time.Time `gorm:"column:published_on; null; type: date"`
	LastModifiedOn  time.Time `gorm:"column:last_modified_on; null; type: date"`
	LastModifiedMsg string    `gorm:"column:last_modified_message; type: character varying(50);"`
	DaysInSale      int32     `gorm:"column:days_in_sale; default: 0"`
	Sold            bool      `gorm:"column:sold; default: false"`
	Promoted        bool      `gorm:"column:promoted; default: false"`
	CreatedOn       time.Time `gorm:"column:created_on; null; type: date; default: now()"`
}
