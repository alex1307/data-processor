package dbmodel

import "time"

type DeletedOnAuditLog struct {
	ID        int32     `gorm:"column:id; primary_key:yes; type: integer; auto_increment:true"`
	ADVERT_ID string    `gorm:"column:advert_id; type: character varying(36); not null"`
	DeletedOn time.Time `gorm:"column:deleted_on; type: date"`
}
