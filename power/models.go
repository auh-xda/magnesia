package power

type PowerInfo struct {
	Vendor   string `json:"vendor"`
	Model    string `json:"model"`
	Serial   string `json:"serial"`
	Status   string `json:"status"`   // Charging / Discharging / Full
	Capacity string `json:"capacity"` // Percentage
}

type WinBattery struct {
	Name                     string
	Manufacturer             string
	SerialNumber             string
	BatteryStatus            uint16
	EstimatedChargeRemaining uint16
}
