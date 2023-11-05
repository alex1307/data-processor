package dbmodel

import (
	"time"
)

type VehicleLog struct {
	ID        uint64    `gorm:"column:id; primary_key:yes; type: bigint; auto_increment:yes"`
	ADVERT_ID string    `gorm:"column:advert_id; type: character varying(36)"`
	Make      string    `gorm:"column:make;not null; type: character varying(25)"`
	Model     string    `gorm:"column:model; not null; type: character varying(50)"`
	Engine    string    `gorm:"column:engine; null; type: character varying(25)"`
	Gearbox   string    `gorm:"column:gearbox; null; type: character varying(25)"`
	Power     int32     `gorm:"column:power; default: 0"`
	Millage   int32     `gorm:"column:millage; default: 0"`
	Year      int32     `gorm:"column:year; default: 0"`
	Currency  string    `gorm:"column:currency; not null; type: character varying(3)"`
	Price     int32     `gorm:"column:price; not null; type: decimal(10,2)"`
	ViewCount int32     `gorm:"column:view_count; default: 0"`
	Equipment int64     `gorm:"column:equipment; type: bigint; default: 0"`
	Phone     string    `gorm:"column:phone; not null; type: character varying(25); default: ''"`
	Location  string    `gorm:"column:location; not null; type: character varying(25); default: ''"`
	Top       bool      `gorm:"column:top; default: false"`
	VIP       bool      `gorm:"column:vip; default: false"`
	Sold      bool      `gorm:"column:sold; default: false"`
	Dealer    bool      `gorm:"column:dealer; default: false"`
	ScrapedOn time.Time `gorm:"column:created_on; not null; type: date"`
	CreatedOn time.Time `gorm:"column:updated_on; null; type: date"`
}
