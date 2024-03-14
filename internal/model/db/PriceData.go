package dbmodel

import "time"

type PriceData struct {
	ID             uint64    `gorm:"column:id; primary_key:yes; type: integer; auto_increment:true"`
	ADVERT_ID      string    `gorm:"column:advert_id; type: character varying(36); not null"`
	Source         string    `gorm:"column:source; default: 0"`
	Currency       string    `gorm:"column:currency;type: character varying(3); not null default: 'EUR'"`
	Price          uint32    `gorm:"column:price; default: 0"`
	EstimatedPrice uint32    `gorm:"column:estimated_price; default: 0"`
	SaveDiff       uint32    `gorm:"column:save_diff; default: 0"`
	OverpricedDiff uint32    `gorm:"column:overpiced_diff; default: 0"`
	Ranges         string    `gorm:"column:ranges; character varying(50);null;"`
	CreatedOn      time.Time `gorm:"column:created_on; null; type: date; default: now()"`
}
