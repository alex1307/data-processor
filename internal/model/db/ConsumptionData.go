package dbmodel

import "time"

type ConsumptionData struct {
	ID              uint64    `gorm:"column:id; primary_key:yes; type: integer; auto_increment:true"`
	ADVERT_ID       string    `gorm:"column:advert_id; type: character varying(36); not null"`
	Source          string    `gorm:"column:source; ; type: character varying(25); not null"`
	Make            string    `gorm:"column:make; type: character varying(50); not null"`
	Model           string    `gorm:"column:model; type: character varying(50); not null"`
	Year            uint32    `gorm:"column:year; default: 0"`
	CO2Emission     uint32    `gorm:"column:co2_emission; default: 0"`
	FuelConsumption float32   `gorm:"column:fuel_consumption; default: 0"`
	KWConsumption   float32   `gorm:"column:kw_consumption; default: 0"`
	CreatedOn       time.Time `gorm:"column:created_on; null; type: date; default: now()"`
}
