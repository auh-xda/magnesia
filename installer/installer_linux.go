package installer

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/auh-xda/magnesia/config"
	"github.com/auh-xda/magnesia/console"
)

func CreateService() error {
	serviceFile := `/etc/systemd/system/magnesia.service`
	serviceContent := `[Unit]
Description=Magnesia Agent
After=network.target

[Service]
ExecStart=/usr/local/bin/magnesia
Restart=always
RestartSec=5
User=root
WorkingDirectory=/usr/local/bin
StandardOutput=append:/var/log/magnesia.log
StandardError=append:/var/log/magnesia.log

[Install]
WantedBy=multi-user.target
`

	// ensure log file exists
	if err := os.MkdirAll("/var/log", 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %v", err)
	}
	if _, err := os.Stat("/var/log/magnesia.log"); os.IsNotExist(err) {
		if err := os.WriteFile("/var/log/magnesia.log", []byte(""), 0644); err != nil {
			return fmt.Errorf("failed to create log file: %v", err)
		}
	}

	// write the service file
	if err := os.WriteFile(serviceFile, []byte(serviceContent), 0644); err != nil {
		return fmt.Errorf("failed to write service file: %v", err)
	}

	// find the current binary
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get current executable path: %v", err)
	}

	// copy binary to /usr/local/bin/magnesia
	if err := exec.Command("cp", exePath, "/usr/local/bin/magnesia").Run(); err != nil {
		return fmt.Errorf("failed to copy binary: %v", err)
	}
	if err := os.Chmod("/usr/local/bin/magnesia", 0755); err != nil {
		return fmt.Errorf("failed to set executable bit: %v", err)
	}

	// reload systemd and enable service
	if err := exec.Command("systemctl", "daemon-reload").Run(); err != nil {
		return fmt.Errorf("daemon-reload failed: %v", err)
	}
	if err := exec.Command("systemctl", "enable", "--now", "magnesia").Run(); err != nil {
		return fmt.Errorf("enable service failed: %v", err)
	}

	console.Success("Magnesia service installed and started successfully (logs: /var/log/magnesia.log)")
	return nil
}

func Uninstall() {
	serviceFile := `/etc/systemd/system/magnesia.service`
	binaryFile := `/usr/local/bin/magnesia`
	logFile := `/var/log/magnesia.log`

	console.Warn("Disabling the service")
	_ = exec.Command("systemctl", "stop", "magnesia").Run()
	_ = exec.Command("systemctl", "disable", "magnesia").Run()

	console.Warn("Removing the daemons")
	_ = os.Remove(serviceFile)
	_ = os.Remove(binaryFile)
	_ = os.Remove(logFile)
	_ = os.RemoveAll(config.Dir())
}
