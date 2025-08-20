//go:build windows
// +build windows

package interceptor

import (
	"fmt"
	"time"

	"github.com/StackExchange/wmi"
	"github.com/auh-xda/magnesia/console"
	"github.com/shirou/gopsutil/v3/cpu"
	"golang.org/x/sys/windows/svc/mgr"
)

func ListServices() ([]WindowsService, error) {

	var services []WindowsService

	console.Info("Getting services for windows")

	// Connect to the Windows service manager
	m, err := mgr.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to service manager: %w", err)
	}
	defer m.Disconnect()

	// List all service names
	serviceNames, err := m.ListServices()
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}

	for _, name := range serviceNames {
		s, err := m.OpenService(name)
		if err != nil {
			continue // skip services we can't open
		}
		defer s.Close()

		// Query config and status
		config, _ := s.Config()
		status, _ := s.Query()

		state := "stopped"
		if status.State == 4 { // SERVICE_RUNNING
			state = "running"
		}

		services = append(services, WindowsService{
			Name:        name,
			DisplayName: config.DisplayName,
			Status:      state,
			StartType:   startTypeToString(config.StartType),
		})
	}

	console.Log(services)

	return services, nil
}

func startTypeToString(t uint32) string {
	switch t {
	case 2:
		return "automatic"
	case 3:
		return "manual"
	case 4:
		return "disabled"
	default:
		return "unknown"
	}
}

func GetInfo() PowerInfo {
	var statuses []BatteryStatus
	_ = wmi.QueryNamespace(wmi.CreateQuery(&statuses, ""), &statuses, "root\\wmi")

	var fullCaps []BatteryFullChargedCapacity
	_ = wmi.QueryNamespace(wmi.CreateQuery(&fullCaps, ""), &fullCaps, "root\\wmi")

	percent := 0
	status := "Unknown"

	if len(statuses) > 0 {
		b := statuses[0]

		// Calculate percentage safely
		if len(fullCaps) > 0 && fullCaps[0].FullChargedCapacity > 0 {
			percent = int((float64(b.RemainingCapacity) / float64(fullCaps[0].FullChargedCapacity)) * 100)
		}

		if b.Charging {
			status = "Charging"
		} else if b.Discharging {
			status = "Discharging"
		} else if b.PowerOnline {
			status = "On AC Power"
		}
	}

	return PowerInfo{
		Status:   status,
		Capacity: fmt.Sprintf("%d%%", percent),
	}
}

func GetCPUInfo() CPUInfo {
	console.Info("Getting CPU info from windows")
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
