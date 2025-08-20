//go:build windows
// +build windows

package power

import (
	"fmt"

	"github.com/StackExchange/wmi"
)

type BatteryStatus struct {
	Charging          bool
	Discharging       bool
	PowerOnline       bool
	RemainingCapacity uint32
	ChargeRate        int32
	DischargeRate     int32
	Voltage           uint32
}

func GetInfo() PowerInfo {
	var statuses []BatteryStatus
	// Query from root\wmi namespace
	q := wmi.CreateQuery(&statuses, "")
	err := wmi.QueryNamespace(q, &statuses, "root\\wmi")
	if err != nil || len(statuses) == 0 {
		return PowerInfo{Status: "No Battery Detected"}
	}

	b := statuses[0]

	status := "Unknown"
	if b.Charging {
		status = "Charging"
	} else if b.Discharging {
		status = "Discharging"
	} else if b.PowerOnline {
		status = "On AC Power"
	}

	return PowerInfo{
		Vendor:   "", // Not in BatteryStatus
		Model:    "", // Not in BatteryStatus
		Serial:   "", // Not in BatteryStatus
		Status:   status,
		Capacity: fmt.Sprintf("%d%%", b.RemainingCapacity),
	}
}
