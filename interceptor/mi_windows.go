//go:build windows
// +build windows

package interceptor

import (
	"fmt"
	"time"

	"github.com/StackExchange/wmi"
	"github.com/auh-xda/magnesia/console"
	"github.com/shirou/gopsutil/v3/cpu"
	"golang.org/x/sys/windows/registry"
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

func GetPowerInfo() (PowerInfo, error) {
	var batteries []win32Battery

	// Query only universally safe fields
	err := wmi.Query("SELECT Name, DeviceID, EstimatedChargeRemaining, BatteryStatus FROM Win32_Battery", &batteries)
	if err != nil {
		return PowerInfo{}, fmt.Errorf("failed to query battery info: %v", err)
	}

	if len(batteries) == 0 {
		return PowerInfo{}, nil
	}

	b := batteries[0]

	status := "Unknown"
	capacity := "Unknown"

	if b.BatteryStatus != nil {
		status = batteryStatusText(*b.BatteryStatus)
	}
	if b.EstimatedChargeRemaining != nil {
		capacity = fmt.Sprintf("%d%%", *b.EstimatedChargeRemaining)
	}

	return PowerInfo{
		Vendor:   "",
		Model:    safeString(b.Name),
		Serial:   safeString(b.DeviceID),
		Status:   status,
		Capacity: capacity,
	}, nil
}

func safeString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func batteryStatusText(code uint16) string {
	switch code {
	case 1:
		return "Discharging"
	case 2:
		return "AC (Charging)"
	case 3:
		return "Fully Charged"
	case 4:
		return "Low"
	case 5:
		return "Critical"
	case 6:
		return "Charging"
	case 7:
		return "Charging and High"
	case 8:
		return "Charging and Low"
	case 9:
		return "Charging and Critical"
	case 10:
		return "Undefined"
	case 11:
		return "Partially Charged"
	default:
		return "Unknown"
	}
}

func GetCPUInfo() (CPUInfo, error) {
	var win32CPUs []win32Processor
	err := wmi.Query("SELECT Manufacturer, Name, NumberOfCores, NumberOfLogicalProcessors, MaxClockSpeed, SocketDesignation FROM Win32_Processor", &win32CPUs)

	// CPU usage stats (works regardless of WMI success/failure)
	usagePerCore, _ := cpu.Percent(time.Second, true)
	overallUsage, _ := cpu.Percent(time.Second, false)

	// If WMI failed → fallback to gopsutil basic info
	if err != nil || len(win32CPUs) == 0 {
		infoStats, errInfo := cpu.Info()
		if errInfo != nil || len(infoStats) == 0 {
			return CPUInfo{}, fmt.Errorf("failed to fetch CPU info via WMI and gopsutil")
		}
		ci := infoStats[0]
		return CPUInfo{
			Manufacturer:      "Unknown",
			Model:             ci.ModelName,
			SpeedMHz:          ci.Mhz,
			TotalCores:        int(ci.Cores),
			LogicalProcessors: int(len(infoStats)), // gopsutil reports logical CPUs
			Sockets:           1,
			CoresPerSocket:    int(ci.Cores),
			Hyperthread:       len(infoStats) > int(ci.Cores),
			UsagePerCore:      usagePerCore,
			OverallUsage:      overallUsage[0],
		}, nil
	}

	// Aggregate multi-socket results
	totalCores := 0
	totalLogical := 0
	sockets := len(win32CPUs)
	manufacturer := win32CPUs[0].Manufacturer
	model := win32CPUs[0].Name
	speed := win32CPUs[0].MaxClockSpeed
	coresPerSocket := int(win32CPUs[0].NumberOfCores)

	for _, c := range win32CPUs {
		totalCores += int(c.NumberOfCores)
		totalLogical += int(c.NumberOfLogicalProcessors)
	}

	info := CPUInfo{
		Manufacturer:      manufacturer,
		Model:             model,
		SpeedMHz:          float64(speed),
		TotalCores:        totalCores,
		LogicalProcessors: totalLogical,
		Sockets:           sockets,
		CoresPerSocket:    coresPerSocket,
		Hyperthread:       totalLogical > totalCores,
		UsagePerCore:      usagePerCore,
		OverallUsage:      overallUsage[0],
	}

	return info, nil
}

func Installations() ([]InstalledSoftware, error) {
	roots := []registry.Key{
		registry.LOCAL_MACHINE,
	}
	paths := []string{
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall`,
		`SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\Uninstall`,
	}

	var software []InstalledSoftware

	for _, root := range roots {
		for _, path := range paths {
			k, err := registry.OpenKey(root, path, registry.READ)
			if err != nil {
				continue
			}

			defer k.Close()

			names, _ := k.ReadSubKeyNames(-1)
			for _, name := range names {
				subk, err := registry.OpenKey(k, name, registry.READ)
				if err != nil {
					continue
				}
				defer subk.Close()

				displayName, _, _ := subk.GetStringValue("DisplayName")
				if displayName == "" {
					continue
				}

				displayVersion, _, _ := subk.GetStringValue("DisplayVersion")
				publisher, _, _ := subk.GetStringValue("Publisher")
				installDate, _, _ := subk.GetStringValue("InstallDate")
				installLocation, _, _ := subk.GetStringValue("InstallLocation")
				uninstallString, _, _ := subk.GetStringValue("UninstallString")
				quietUninstall, _, _ := subk.GetStringValue("QuietUninstallString")
				iconPath, _, _ := subk.GetStringValue("DisplayIcon")
				helpLink, _, _ := subk.GetStringValue("HelpLink")
				infoURL, _, _ := subk.GetStringValue("URLInfoAbout")
				installSource, _, _ := subk.GetStringValue("InstallSource")

				size, _, _ := subk.GetIntegerValue("EstimatedSize") // in KB

				software = append(software, InstalledSoftware{
					Name:            displayName,
					Version:         displayVersion,
					Vendor:          publisher,
					InstallDate:     installDate,
					InstallLocation: installLocation,
					UninstallString: uninstallString,
					QuietUninstall:  quietUninstall,
					EstimatedSize:   formatSizeKB(size),
					IconPath:        iconPath,
					HelpLink:        helpLink,
					InfoURL:         infoURL,
					InstallSource:   installSource,
				})

				software = append(software, InstalledSoftware{
					Name:        displayName,
					Version:     displayVersion,
					Vendor:      publisher,
					InstallDate: installDate,
				})
			}
		}
	}
	return software, nil
}

func formatSizeKB(sizeKB uint64) string {
	if sizeKB == 0 {
		return ""
	}

	if sizeKB < 1024 {
		return fmt.Sprintf("%d KB", sizeKB)
	}

	sizeMB := float64(sizeKB) / 1024
	if sizeMB < 1024 {
		return fmt.Sprintf("%.1f MB", sizeMB)
	}

	sizeGB := sizeMB / 1024
	return fmt.Sprintf("%.1f GB", sizeGB)
}
