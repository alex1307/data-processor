package dbmodel

type IDData struct {
	ID        uint64 `gorm:"column:id; primary_key:yes; type: integer; auto_increment:true"`
	ADVERT_ID string `gorm:"column:advert_id; type: character varying(50); not null"`
	Source    string `gorm:"column:source; ; type: character varying(25); not null"`
	CreatedOn string `gorm:"column:created_on; null; type: date; default: now()"`
}
