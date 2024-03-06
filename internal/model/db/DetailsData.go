package dbmodel

import "time"

type DetailsData struct {
	ID         uint64    `gorm:"column:id; primary_key:yes; type: integer; auto_increment:true"`
	ADVERT_ID  string    `gorm:"column:advert_id; type: character varying(36); not null"`
	Source     string    `gorm:"column:source; type: character varying(25); not null"`
	Phone      string    `gorm:"column:phone; type: character varying(50);null;"`
	Location   string    `gorm:"column:location; type: character varying(100);null;"`
	ViewCount  uint32    `gorm:"column:view_count; default: 0;null;"`
	Equipment  uint64    `gorm:"column:equipment; default: 0;null;"`
	SellerName string    `gorm:"column:seller_name; type: character varying(50);null;"`
	IsDealer   bool      `gorm:"column:is_dealer; default: false;null;"`
	CreatedOn  time.Time `gorm:"column:created_on; null; type: date; default: now()"`
}
