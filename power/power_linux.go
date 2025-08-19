//go:build linux
// +build linux

package power

import (
	"os"
	"strings"
)

func GetInfo() PowerInfo {
	return PowerInfo{
		Vendor:   readFile("/sys/class/power_supply/BAT0/manufacturer"),
		Model:    readFile("/sys/class/power_supply/BAT0/model_name"),
		Serial:   readFile("/sys/class/power_supply/BAT0/serial_number"),
		Status:   readFile("/sys/class/power_supply/BAT0/status"),
		Capacity: readFile("/sys/class/power_supply/BAT0/capacity"),
	}
}

func readFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return "N/A"
	}
	return strings.TrimSpace(string(data))
}
