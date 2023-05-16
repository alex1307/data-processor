package dbmodel

import (
	"time"
)

type Equipment struct {
	ID                    string    `gorm:"column:id; primary_key:yes; type: character varying(36)"`
	Parktronic            bool      `gorm:"column:parktronic; default: false"`
	AdaptiveCruiseControl bool      `gorm:"column:adaptive_cruise_control; default: false"`
	DVD                   bool      `gorm:"column:dvd; default: false"`
	Insurance             bool      `gorm:"column:insurance; default: false"`
	HeatedSeats           bool      `gorm:"column:heated_seats; default: false"`
	AdapterAirSuspension  bool      `gorm:"column:adapter_air_suspension; default: false"`
	Registration          bool      `gorm:"column:registration; default: false"`
	NewImport             bool      `gorm:"column:new_import; default: false"`
	LeatherSeats          bool      `gorm:"column:leather_seats; default: false"`
	Sunroof               bool      `gorm:"column:sunroof; default: false"`
	Methane               bool      `gorm:"column:methane; default: false"`
	Leasing               bool      `gorm:"column:leasing; default: false"`
	FourWheelDrive        bool      `gorm:"column:four_wheel_drive; default: false"`
	LPG                   bool      `gorm:"column:lpg; default: false"`
	CruiseControl         bool      `gorm:"column:cruise_control; default: false"`
	UpdatedOn             time.Time `gorm:"column:updated_on;"`
}
