package main

import (
	"github.com/auh-xda/magnesia/interceptor"
)

type Config struct {
	Version  string `json:"version"`
	UUID     string `json:"uuid"`
	Momentum string `json:"server"`
	Interval string `json:"interval"`
	Channel  string `json:"channel"`
	ClientID string `json:"client_id"`
}

type AuthRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	AuthToken    string `json:"auth_token"`
	ApiKey       string `json:"api_key"`
}

type AuthResponse struct {
	Config  Config `json:"config"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type Magnesia struct {
	Action       string `json:"action"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	AuthToken    string `json:"auth_token"`
	ApiKey       string `json:"api_key"`
}

type Websocket struct {
	MagnesiaUid     string `json:"magnesia_uuid"`
	MagnesiaPayload any    `json:"magnesia_payload"`
	MagnesiaChannel string `json:"magnesia_channel"`
	MagnesiaSiteId  string `json:"magnesia_site_id"`
}

type Interface struct {
	Name        string   `json:"name"`
	MacAddress  string   `json:"mac"`
	IPAddresses []string `json:"ip_addresses"`
}

type Intercept struct {
	Version        string                `json:"version"`
	SerialNumber   string                `json:"product_serial"`
	Hostname       string                `json:"hostname"`
	PublicIP       string                `json:"public_ip"`
	OS             string                `json:"os"`
	OSVersion      string                `json:"os_version"`
	UpTime         uint64                `json:"uptime"`
	BootTime       uint64                `json:"boot_time"`
	HostID         string                `json:"host_id"`
	PlatformFamily string                `json:"family"`
	Interfaces     []Interface           `json:"interfaces"`
	Power          interceptor.PowerInfo `json:"power"`
	Memory         MemoryInfo            `json:"memory"`
	DiskInfo       []DiskInfo            `json:"disks"`
	CPUInfo        CPUInfo               `json:"cpu"`
}

type CPUInfo struct {
	Manufacturer      string    `json:"manufacturer"`
	SpeedMHz          float64   `json:"cpu_speed_mhz"`
	TotalCores        int       `json:"cores"`
	Model             string    `json:"model"`
	Sockets           int       `json:"sockets"`
	CoresPerSocket    int       `json:"cores_per_socket"`
	LogicalProcessors int       `json:"logical_processors"`
	Hyperthread       bool      `json:"hyperthread"`
	UsagePerCore      []float64 `json:"usage_per_core"`
	OverallUsage      float64   `json:"overall_usage"`
}

type MemoryInfo struct {
	Total uint64  `json:"total"`
	Used  uint64  `json:"used"`
	Free  uint64  `json:"free"`
	Usage float64 `json:"usage_percent"`
}

type DiskInfo struct {
	Device       string  `json:"device"`
	MountPoint   string  `json:"mountpoint"`
	Fstype       string  `json:"fstype"`
	Total        uint64  `json:"total"`
	Used         uint64  `json:"used"`
	Free         uint64  `json:"free"`
	UsagePercent float64 `json:"usage_percent"`
}

type ProcessInfo struct {
	PID        int32   `json:"pid"`
	PPID       int32   `json:"ppid"`
	Name       string  `json:"name"`
	Exe        string  `json:"exe"`
	Cmdline    string  `json:"cmdline"`
	Username   string  `json:"username"`
	Status     string  `json:"status"`
	CPUPercent float64 `json:"cpu_percent"`
	MemoryMB   float32 `json:"memory_mb"`
	CreateTime int64   `json:"create_time"`
	NumThreads int32   `json:"num_threads"`
	Nice       int32   `json:"nice"`
}
