package dbmodel

import "time"

type BasicData struct {
	ID        uint64    `gorm:"column:id; primary_key:yes; type: integer; auto_increment:true"`
	ADVERT_ID string    `gorm:"column:advert_id; type: character varying(36); not null"`
	Source    string    `gorm:"column:source; type: character varying(25); not null"`
	Title     string    `gorm:"column:title; type: character varying(100); not null"`
	Make      string    `gorm:"column:make; type: character varying(50); not null"`
	Model     string    `gorm:"column:model; type: character varying(50); not null"`
	Year      uint32    `gorm:"column:year; default: 0"`
	Currency  string    `gorm:"column: currency;type: character varying(3); not null default: 'EUR'"`
	Price     uint32    `gorm:"column:price; default: 0"`
	Millage   uint32    `gorm:"column:millage; default: 0"`
	Month     uint32    `gorm:"column:month; default: 0"`
	Engine    string    `gorm:"column:engine; type: character varying(50); not null"`
	Gearbox   string    `gorm:"column:gearbox; type: character varying(50); not null"`
	CC        uint32    `gorm:"column:cc; default: 0"`
	PowerPS   uint32    `gorm:"column:power_ps; default: 0"`
	PowerKW   uint32    `gorm:"column:power_kw; default: 0"`
	CreatedOn time.Time `gorm:"column:created_on; null; type: date; default: now()"`
}
