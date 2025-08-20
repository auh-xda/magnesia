//go:build darwin
// +build darwin

package interceptor

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/auh-xda/magnesia/console"
)

func ListServices() ([]DarwinService, error) {
	var services []DarwinService

	console.Info("Getting services for darwin/mac")

	cmd := exec.Command("launchctl", "list")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run launchctl: %w", err)
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))
	firstLine := true
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		if firstLine { // skip header
			firstLine = false
			continue
		}
		if len(line) < 3 {
			continue
		}

		// PID
		pid := 0
		if line[0] != "-" {
			pid, _ = strconv.Atoi(line[0])
		}

		// Status (0 means running ok, nonzero means error)
		status := "stopped"
		if pid > 0 {
			status = "running"
		}

		services = append(services, DarwinService{
			Label:  line[2],
			Status: status,
			PID:    pid,
		})
	}

	return services, nil
}

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

func GetCPUInfo() CPUInfo {
	listOfCpus, err := cpu.Info()
	if err != nil || len(listOfCpus) == 0 {
		return CPUInfo{}
	}

	uniqueCores := make(map[string]struct{})
	uniqueSockets := make(map[string]struct{})

	for _, c := range listOfCpus {
		uniqueCores[c.CoreID] = struct{}{}
		uniqueSockets[c.PhysicalID] = struct{}{}
	}

	totalCores := len(uniqueCores)
	totalSockets := len(uniqueSockets)
	if totalSockets == 0 {
		totalSockets = 1 // avoid divide by zero
	}
	coresPerSocket := totalCores / totalSockets

	logicalProcs := len(listOfCpus)

	// CPU usage percentages
	usagePercents, _ := cpu.Percent(1*time.Second, false)        // overall
	usagePercentsCoreWise, _ := cpu.Percent(1*time.Second, true) // per core

	overallUsage := 0.0

	if len(usagePercents) > 0 {
		overallUsage = usagePercents[0]
	}

	cpuInfo := CPUInfo{
		Manufacturer:      listOfCpus[0].VendorID,
		Model:             listOfCpus[0].ModelName,
		SpeedMHz:          listOfCpus[0].Mhz,
		TotalCores:        totalCores,
		Sockets:           totalSockets,
		CoresPerSocket:    coresPerSocket,
		LogicalProcessors: logicalProcs,
		Hyperthread:       logicalProcs > totalCores,
		UsagePerCore:      usagePercentsCoreWise,
		OverallUsage:      overallUsage,
	}

	return cpuInfo
}
