package jsonmodel

type EquipmentDTO struct {
	Parktronic            bool `json:"parktronic"`
	AdaptiveCruiseControl bool `json:"adaptive_cruise_control"`
	DVD                   bool `json:"dvd"`
	Insurance             bool `json:"insurance"`
	HeatedSeats           bool `json:"heated_seats"`
	AdapterAirSuspension  bool `json:"adapter_air_suspension"`
	Registration          bool `json:"registration"`
	FullyServiced         bool `json:"fully_serviced"`
	NewImport             bool `json:"new_import"`
	LeatherSeats          bool `json:"leather_seats"`
	Sunroof               bool `json:"sunroof"`
	Methane               bool `json:"methane"`
	Leasing               bool `json:"leasing"`
	FourWheelDrive        bool `json:"four_wheel_drive"`
	LPG                   bool `json:"lpg"`
	CruiseControl         bool `json:"cruise_control"`
}
