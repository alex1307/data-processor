package dbmodel

import (
	"time"
)

type Equipment struct {
	ID                    int64     `gorm:"column:id; primary_key:yes;autoIncrement:false;"`
	Parktronic            bool      `gorm:"column:parktronic; default: false"`
	AdaptiveCruiseControl bool      `gorm:"column:adaptive_cruise_control; default: false"`
	DVD                   bool      `gorm:"column:dvd; default: false"`
	Insurance             bool      `gorm:"column:insurance; default: false"`
	HeatedSeats           bool      `gorm:"column:heated_seats; default: false"`
	AdapterAirSuspension  bool      `gorm:"column:adapter_air_suspension; default: false"`
	Registration          bool      `gorm:"column:registration; default: false"`
	FullyServiced         bool      `gorm:"column:fully_serviced; default: false"`
	NewImport             bool      `gorm:"column:new_import; default: false"`
	LeatherSeats          bool      `gorm:"column:leather_seats; default: false"`
	Sunroof               bool      `gorm:"column:sunroof; default: false"`
	Methane               bool      `gorm:"column:methane; default: false"`
	Leasing               bool      `gorm:"column:leasing; default: false"`
	FourWheelDrive        bool      `gorm:"column:four_wheel_drive; default: false"`
	LPG                   bool      `gorm:"column:lpg; default: false"`
	CruiseControl         bool      `gorm:"column:cruise_control; default: false"`
	CreatedOn             time.Time `gorm:"column:created_on; not null; type: date"`
}
