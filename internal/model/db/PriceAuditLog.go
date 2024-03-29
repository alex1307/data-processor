package dbmodel

import "time"

type PriceAuditLog struct {
	ID        int32     `gorm:"column:id; primary_key:yes; type: integer; auto_increment:true"`
	ADVERT_ID string    `gorm:"column:advert_id; type: character varying(36); not null"`
	NewValue  int32     `gorm:"column:new_value; default: 0"`
	OldValue  int32     `gorm:"column:old_value; default: 0"`
	DiffValue int32     `gorm:"column:difference; default: 0"`
	UpdatedOn time.Time `gorm:"column:updated_on; null; type: date"`
}
