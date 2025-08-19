package main

import (
	"bytes"
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strings"
	"time"

	console "github.com/auh-xda/magnesia/helpers/console"
	"github.com/auh-xda/magnesia/power"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

func (magnesia Magnesia) Intercept() {
	console.Info("Collecting Information ...")

	start := time.Now()

	intercept := Intercept{}

	intercept.Version = version
	intercept.SerialNumber = getProductSerial()

	publicIP, err := exec.Command("curl", "-s", "https://ifconfig.me").Output()
	if err == nil {
		intercept.PublicIP = strings.TrimSpace(string(publicIP))
	}

	info, err := host.Info()

	if err == nil {
		intercept.OS = info.OS
		intercept.OSVersion = info.PlatformVersion
		intercept.Hostname = info.Hostname
		intercept.UpTime = info.Uptime
		intercept.BootTime = info.BootTime
		intercept.PlatformFamily = info.PlatformFamily
		intercept.HostID = info.HostID
	}

	intercept.Power = power.GetInfo()
	intercept.Interfaces = getDeviceInterfaces()
	intercept.Memory = getMemoryInfo()
	intercept.DiskInfo = getDiskInfo()
	intercept.CPUInfo = getCPUInfo()

	console.Log(intercept)

	time := time.Since(start).Seconds()
	console.Success(fmt.Sprintf("Information pulled up in %0.2f s", time))

	Websocket{MagnesiaPayload: intercept}.SendData()
}

func getProductSerial() string {
	switch runtime.GOOS {
	case "windows":
		out, err := exec.Command("wmic", "bios", "get", "serialnumber").Output()
		if err != nil {
			return "--"
		}
		lines := strings.Split(string(out), "\n")
		if len(lines) > 1 {
			return strings.TrimSpace(lines[1])
		}
		return "--"

	case "linux":
		// Try reading from DMI
		out, err := exec.Command("cat", "/sys/class/dmi/id/product_serial").Output()
		if err == nil {
			return strings.TrimSpace(string(out))
		}
		// Fallback to dmidecode
		out, err = exec.Command("dmidecode", "-s", "system-serial-number").Output()
		if err == nil {
			return strings.TrimSpace(string(out))
		}
		return "--"

	case "darwin":
		out, err := exec.Command("system_profiler", "SPHardwareDataType").Output()
		if err != nil {
			return "--"
		}
		lines := bytes.Split(out, []byte("\n"))
		for _, l := range lines {
			line := strings.TrimSpace(string(l))
			if strings.HasPrefix(line, "Serial Number") {
				parts := strings.Split(line, ":")
				if len(parts) == 2 {
					return strings.TrimSpace(parts[1])
				}
			}
		}
		return "--"
	}

	return "--"
}

func getCPUInfo() CPUInfo {
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

func getDeviceInterfaces() []Interface {

	ifaces, _ := net.Interfaces()

	interfaces := make([]Interface, 0)

	for _, iface := range ifaces {

		// Interfaces to Skip :
		// Skip down/loopback interfaces
		// virtual interfaces

		if iface.Flags&net.FlagUp == 0 ||
			iface.Flags&net.FlagLoopback != 0 ||
			strings.HasPrefix(iface.Name, "docker") ||
			strings.HasPrefix(iface.Name, "br-") ||
			strings.HasPrefix(iface.Name, "veth") {
			continue
		}

		systemInterface := Interface{}

		systemInterface.Name = iface.Name

		if iface.HardwareAddr != nil {
			systemInterface.MacAddress = iface.HardwareAddr.String()
		}

		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
				if ipNet.IP.To4() != nil {
					systemInterface.IPAddresses = append(systemInterface.IPAddresses, ipNet.IP.String())
				}
			}
		}

		interfaces = append(interfaces, systemInterface)
	}

	return interfaces
}

func getMemoryInfo() MemoryInfo {
	vm, error := mem.VirtualMemory()

	if error != nil {
		console.Log("Failed to get Memory details")
	}

	return MemoryInfo{
		Total: vm.Total,
		Used:  vm.Used,
		Free:  vm.Free,
		Usage: vm.UsedPercent,
	}
}

func getDiskInfo() []DiskInfo {
	partitions, _ := disk.Partitions(false)
	var disks []DiskInfo

	for _, p := range partitions {
		if p.Fstype == "tmpfs" || p.Fstype == "squashfs" || p.Device == "" {
			continue // skip ephemeral/virtual mounts
		}

		usage, err := disk.Usage(p.Mountpoint)

		if err == nil {
			disks = append(disks, DiskInfo{
				MountPoint:   p.Mountpoint,
				Total:        usage.Total,
				Used:         usage.Used,
				Free:         usage.Free,
				UsagePercent: usage.UsedPercent,
				Device:       p.Device,
				Fstype:       p.Fstype,
			})
		}
	}
	return disks
}
