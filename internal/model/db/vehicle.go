package dbmodel

import (
	"time"
)

type Vehicle struct {
	ID        string    `gorm:"column:id; primary_key:yes; type: character varying(36)"`
	Make      string    `gorm:"column:make;not null; type: character varying(25)"`
	Model     string    `gorm:"column:model; not null; type: character varying(50)"`
	Engine    string    `gorm:"column:engine; null; type: character varying(25)"`
	Gearbox   string    `gorm:"column:gearbox; null; type: character varying(25)"`
	Power     int       `gorm:"column:power; default: 0"`
	Millage   int32     `gorm:"column:millage; default: 0"`
	Year      int       `gorm:"column:year; default: 0"`
	Currency  string    `gorm:"column:currency; not null; type: character varying(3)"`
	Price     int32     `gorm:"column:price; not null; type: decimal(10,2)"`
	ViewCount int       `gorm:"column:view_count; default: 0"`
	Equipment int64     `gorm:"column:equipment; type: bigint; default: 0"`
	Source    string    `gorm:"column:source; not null; type: character varying(25); default: 'N/A'"`
	Phone     string    `gorm:"column:phone; not null; type: character varying(25); default: 'N/A'"`
	Promoted  bool      `gorm:"column:promoted; default: false"`
	Sold      bool      `gorm:"column:sold; default: false"`
	CreatedOn time.Time `gorm:"column:created_on; not null; type: date"`
	UpdatedOn time.Time `gorm:"column:updated_on; null; type: date"`
}
