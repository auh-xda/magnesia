// *******************************************************************
// *     __  __                                        _             *
// *    |  \/  |   __ _    __ _   _ __     ___   ___  (_)   __ _     *
// *    | |\/| |  / _` |  / _` | | '_ \   / _ \ / __| | |  / _` |    *
// *    | |  | | | (_| | | (_| | | | | | |  __/ \__ \ | | | (_| |    *
// *    |_|  |_|  \__,_|  \__, | |_| |_|  \___| |___/ |_|  \__,_|    *
// *                      |___/                                      *
// *******************************************************************

package interceptor

type BatteryStatus struct {
	Charging          bool
	Discharging       bool
	PowerOnline       bool
	RemainingCapacity uint32
	ChargeRate        int32
	DischargeRate     int32
	Voltage           uint32
}

type LinuxService struct {
	Name        string `json:"name"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

type WindowsService struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Status      string `json:"status"`
	StartType   string `json:"start_type"`
}

type DarwinService struct {
	Label  string `json:"label"`
	Status string `json:"status"`
	PID    int    `json:"pid,omitempty"`
}

type PowerInfo struct {
	Vendor   string `json:"vendor"`
	Model    string `json:"model"`
	Serial   string `json:"serial"`
	Status   string `json:"status"`
	Capacity string `json:"capacity"`
}

type WinBattery struct {
	Name                     string
	Manufacturer             string
	SerialNumber             string
	BatteryStatus            uint16
	EstimatedChargeRemaining uint16
}
