package dbmodel

import "time"

type Status struct {
	ID        string    `gorm:"column:id; primary_key:yes; type: character varying(36)"`
	Status    string    `gorm:"column:status; not null; type: character varying(25)"`
	UpdatedOn time.Time `gorm:"column:updated_on; type: date"`
}
