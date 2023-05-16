package dbmodel

import "time"

type LogData struct {
	ID          string    `gorm:"column:id; primary_key:yes; type: character varying(36)"`
	Filename    string    `gorm:"column:filename; not null; type: character varying(255)"`
	Status      string    `gorm:"column:status; not null; type: character varying(25)"`
	LastRunDate time.Time `gorm:"column:last_success_date; type: date"`
}
