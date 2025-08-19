package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	client "github.com/auh-xda/magnesia/helpers/client"
	console "github.com/auh-xda/magnesia/helpers/console"
	"github.com/common-nighthawk/go-figure"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

func (magnesia Magnesia) Install() {
	console.SetColor("yellow")
	myFigure := figure.NewFigure("Magnesia", "", true)
	myFigure.Print()
	console.ResetColor()

	console.Info("Installing...")
	config, error := authenticateServer(magnesia)

	if error != nil {
		console.Error("Autentication Failed")
		return
	}

	if error := createConfigFile(config); error != nil {
		console.Error(error.Error())
		return
	}

	magnesia.Intercept()
}

func createConfigFile(config Config) error {
	configDir := "/magnesia"

	console.Info("Generating Magnesia configurations")

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	configFile := filepath.Join(configDir, "config.json")

	jsonData, err := json.MarshalIndent(config, "", "  ")

	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	if err := os.WriteFile(configFile, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	console.Success("Config file updated")

	return nil
}

func authenticateServer(Agent Magnesia) (Config, error) {

	config := Config{}

	authPayload := AuthRequest{
		AuthToken:    Agent.AuthToken,
		ClientID:     Agent.ClientID,
		ClientSecret: Agent.ClientSecret,
		ApiKey:       Agent.ApiKey,
	}

	console.Info("Authenticating with server...")

	response, err := client.Post(authEndpoint, authPayload)

	if err != nil {
		return config, err
	}

	var Auth AuthResponse

	err = json.Unmarshal(response.Body(), &Auth)

	if err != nil {
		return config, fmt.Errorf("error unmarshaling JSON: %s", err)
	}

	if !Auth.Success {
		return config, fmt.Errorf("%s", Auth.Message)
	}

	console.Success("Authentication Sucessfull")

	return Auth.Config, nil
}

func (magnesia Magnesia) Intercept() {
	console.Info("Collecting Information ...")

	start := time.Now()

	intercept := Intercept{}

	intercept.Version = version

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

	intercept.Interfaces = getDeviceInterfaces()
	intercept.Memory = getMemoryInfo()
	intercept.DiskInfo = getDiskInfo()
	intercept.CPUInfo = getCPUInfo()

	console.Log(intercept)

	time := time.Since(start).Seconds()
	console.Success(fmt.Sprintf("Information pulled up in %0.2f s", time))

	Websocket{MagnesiaPayload: intercept}.SendData()
}

func getCPUInfo() CPUInfo {

	listOfCpus, _ := cpu.Info()

	uniqueCores := make(map[string]struct{})
	uniqueSockets := make(map[string]struct{})

	for _, c := range listOfCpus {
		uniqueCores[c.CoreID] = struct{}{}
		uniqueSockets[c.PhysicalID] = struct{}{}
	}

	totalCores := len(uniqueCores)
	totalSockets := len(uniqueSockets)
	logicalProcs := len(listOfCpus)

	usagePercents, _ := cpu.Percent(1*time.Second, false)
	usagePercentsCoreWise, _ := cpu.Percent(1*time.Second, true)

	return CPUInfo{
		Manufacturer:      listOfCpus[0].VendorID,
		SpeedMHz:          listOfCpus[0].Mhz,
		TotalCores:        totalCores,
		Model:             listOfCpus[0].ModelName,
		Sockets:           totalSockets,
		CoresPerSocket:    totalCores / totalSockets,
		LogicalProcessors: logicalProcs,
		Hyperthread:       logicalProcs > totalCores,
		OverallUsage:      usagePercents[0],
		UsagePerCore:      usagePercentsCoreWise,
	}
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
