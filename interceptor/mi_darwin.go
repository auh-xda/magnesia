//go:build darwin
// +build darwin

package interceptor

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/auh-xda/magnesia/console"
	"github.com/shirou/gopsutil/v3/cpu"
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

func GetPowerInfo() (PowerInfo, error) {
	pi := PowerInfo{}
	out, err := exec.Command("ioreg", "-rc", "AppleSmartBattery").Output()
	if err != nil {
		return PowerInfo{}, err
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
	return pi, nil
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

func GetCPUInfo() (CPUInfo, error) {
	listOfCpus, err := cpu.Info()
	if err != nil || len(listOfCpus) == 0 {
		return CPUInfo{}, err
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

	return cpuInfo, nil
}

func Installations() ([]InstalledSoftware, error) {
	cmd := exec.Command("system_profiler", "SPApplicationsDataType", "-json")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	var data SystemProfiler
	if err := json.Unmarshal(out.Bytes(), &data); err != nil {
		return nil, err
	}

	var installed []InstalledSoftware

	for _, app := range data.Applications {
		installed = append(installed, InstalledSoftware{
			Name:            app.Name,
			Version:         app.Version,
			Vendor:          app.ObtainedFrom,
			InstallLocation: app.Path,
			InstallSource:   "system_profiler",
		})
	}

	return installed, nil
}
