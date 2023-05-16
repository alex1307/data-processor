package dbmodel

import (
	"time"
)

type Details struct {
	ID        string    `gorm:"column:id; primary_key:yes; type: character varying(36)"`
	Currency  string    `gorm:"column:currency; not null; type: character varying(3)"`
	Price     int32     `gorm:"column:price; not null; type: decimal(10,2)"`
	ViewCount int       `gorm:"column:view_count; default: 0"`
	Promoted  bool      `gorm:"column:promoted; default: false"`
	Sold      bool      `gorm:"column:sold; default: false"`
	Equipment int64     `gorm:"column:equipment; type: bigint; default: 0"`
	CreatedOn time.Time `gorm:"column:created_on; default: now(); not null; type: date"`
	UpdatedOn time.Time `gorm:"column:updated_on; null"`
}
