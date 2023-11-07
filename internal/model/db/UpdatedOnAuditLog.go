package dbmodel

import "time"

type UpdatedOnAuditLog struct {
	ID        int32     `gorm:"column:id; primary_key:yes; type: integer; auto_increment:true"`
	ADVERT_ID string    `gorm:"column:advert_id; type: character varying(36); not null"`
	NewValue  time.Time `gorm:"column:new_value; type: date"`
	OldValue  time.Time `gorm:"column:old_value; type: date"`
	UpdatedOn time.Time `gorm:"column:updated_on; null; type: date; default: now()"`
}
