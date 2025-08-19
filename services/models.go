package services

// linux.go
type LinuxService struct {
	Name        string `json:"name"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

// windows.go
type WindowsService struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Status      string `json:"status"`
	StartType   string `json:"start_type"`
}

// darwin.go
type DarwinService struct {
	Label  string `json:"label"`
	Status string `json:"status"`
	PID    int    `json:"pid,omitempty"`
}
