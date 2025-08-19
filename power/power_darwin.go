//go:build darwin
// +build darwin

package power

import (
	"os/exec"
	"strconv"
	"strings"
)

func GetInfo() PowerInfo {
	pi := PowerInfo{}
	out, err := exec.Command("ioreg", "-rc", "AppleSmartBattery").Output()
	if err != nil {
		return PowerInfo{}
	}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, `"Manufacturer"`) {
			pi.Vendor = extractStringValue(line)
		} else if strings.Contains(line, `"DeviceName"`) {
			pi.Model = extractStringValue(line)
		} else if strings.Contains(line, `"SerialNumber"`) {
			pi.Serial = extractStringValue(line)
		} else if strings.Contains(line, `"IsCharging"`) {
			if strings.Contains(line, "Yes") || strings.Contains(line, "true") {
				pi.Status = "Charging"
			} else {
				pi.Status = "Discharging"
			}
		} else if strings.Contains(line, `"CurrentCapacity"`) {
			current := extractIntValue(line)
			max := getDarwinMaxCapacity(lines)
			if max > 0 {
				pi.Capacity = strconv.Itoa((current*100)/max) + "%"
			}
		}
	}
	return pi
}

// Helpers to parse ioreg values
func extractStringValue(line string) string {
	parts := strings.Split(line, "=")
	if len(parts) < 2 {
		return ""
	}
	return strings.Trim(strings.TrimSpace(parts[1]), `"`)
}

func extractIntValue(line string) int {
	parts := strings.Split(line, "=")
	if len(parts) < 2 {
		return 0
	}
	val, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
	return val
}

func getDarwinMaxCapacity(lines []string) int {
	for _, line := range lines {
		if strings.Contains(line, `"MaxCapacity"`) {
			return extractIntValue(line)
		}
	}
	return 0
}
