package interceptor

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

type LinuxService struct {
	Name        string `json:"name"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

type WindowsService struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Status      string `json:"status"`
	StartType   string `json:"start_type"`
}

type BatteryFullChargedCapacity struct {
	FullChargedCapacity uint32
}

type DarwinService struct {
	Label  string `json:"label"`
	Status string `json:"status"`
	PID    int    `json:"pid,omitempty"`
}

type PowerInfo struct {
	Vendor   string `json:"vendor,omitempty"`
	Model    string `json:"model,omitempty"`
	Serial   string `json:"serial,omitempty"`
	Status   string `json:"status"`
	Capacity string `json:"capacity"`
}

type BatteryStatus struct {
	Charging          bool
	Discharging       bool
	PowerOnline       bool
	RemainingCapacity uint32
}

type win32Processor struct {
	Manufacturer              string
	Name                      string
	NumberOfCores             uint32
	NumberOfLogicalProcessors uint32
	MaxClockSpeed             uint32
	SocketDesignation         string
}

type win32Battery struct {
	Name                     *string
	DeviceID                 *string
	EstimatedChargeRemaining *uint16
	BatteryStatus            *uint16
}

type InstalledSoftware struct {
	Name            string `json:"name"`
	Version         string `json:"version,omitempty"`
	Vendor          string `json:"vendor,omitempty"`
	InstallDate     string `json:"install_date,omitempty"`
	InstallLocation string `json:"install_location,omitempty"`
	UninstallString string `json:"uninstall_string,omitempty"`
	QuietUninstall  string `json:"quiet_uninstall,omitempty"`
	EstimatedSize   string `json:"estimated_size,omitempty"`
	IconPath        string `json:"icon_path,omitempty"`
	HelpLink        string `json:"help_link,omitempty"`
	InfoURL         string `json:"info_url,omitempty"`
	InstallSource   string `json:"install_source,omitempty"`
}

type SystemProfiler struct {
	Applications []Application `json:"SPApplicationsDataType"`
}

type Application struct {
	Name         string `json:"_name"`
	Version      string `json:"version"`
	Path         string `json:"path"`
	ObtainedFrom string `json:"obtained_from,omitempty"`
}
