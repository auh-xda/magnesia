//go:build windows
// +build windows

package power

import (
	"fmt"

	"github.com/StackExchange/wmi"
)

func GetInfo() PowerInfo {
	var dst []WinBattery
	err := wmi.Query("SELECT * FROM Win32_Battery", &dst)
	if err != nil || len(dst) == 0 {
		return PowerInfo{}
	}
	b := dst[0]
	status := "Unknown"
	switch b.BatteryStatus {
	case 1:
		status = "Discharging"
	case 2:
		status = "Charging"
	case 3:
		status = "Fully Charged"
	}

	return PowerInfo{
		Vendor:   b.Manufacturer,
		Model:    b.Name,
		Serial:   b.SerialNumber,
		Status:   status,
		Capacity: fmt.Sprintf("%d%%", b.EstimatedChargeRemaining),
	}
}
