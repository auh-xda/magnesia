//go:build linux
// +build linux

package interceptor

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/auh-xda/magnesia/console"
	"github.com/shirou/gopsutil/v3/cpu"
)

func ListServices() ([]LinuxService, error) {
	var services []LinuxService

	console.Info("Getting services for Linux")

	// Run systemctl to get list of services
	cmd := exec.Command("systemctl",
		"list-units",
		"--type=service",
		"--all",
		"--no-legend",
		"--no-pager",
		"--plain")

	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run systemctl: %w", err)
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		if len(line) < 4 {
			continue
		}

		unit := line[0] // e.g. ssh.service

		active := line[2] // "active"/"inactive"/"failed"
		description := strings.Join(line[4:], " ")

		state := "stopped"
		if active == "active" {
			state = "running"
		}

		services = append(services, LinuxService{
			Name:        strings.TrimSuffix(unit, ".service"),
			Description: description,
			Status:      state,
		})
	}

	return services, nil
}

func GetPowerInfo() (PowerInfo, error) {
	return PowerInfo{
		Vendor:   readFile("/sys/class/power_supply/BAT0/manufacturer"),
		Model:    readFile("/sys/class/power_supply/BAT0/model_name"),
		Serial:   readFile("/sys/class/power_supply/BAT0/serial_number"),
		Status:   readFile("/sys/class/power_supply/BAT0/status"),
		Capacity: readFile("/sys/class/power_supply/BAT0/capacity"),
	}, nil
}

func readFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return "N/A"
	}
	return strings.TrimSpace(string(data))
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
	return []InstalledSoftware{}, nil
}
