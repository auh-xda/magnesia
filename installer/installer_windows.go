//go:build windows
// +build windows

package installer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/auh-xda/magnesia/config"
	"github.com/auh-xda/magnesia/console"
)

func CreateService() error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %v", err)
	}

	targetBin := `C:\Program Files\Magnesia\magnesia.exe`
	if err := os.MkdirAll(filepath.Dir(targetBin), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}
	if err := exec.Command("copy", exePath, targetBin).Run(); err != nil {
		return fmt.Errorf("failed to copy binary: %v", err)
	}

	// Create Windows service
	serviceName := "MagnesiaAgent"
	createCmd := exec.Command("sc.exe", "create", serviceName,
		"binPath=", targetBin,
		"start=", "auto")
	if err := createCmd.Run(); err != nil {
		return fmt.Errorf("failed to create service: %v", err)
	}

	// Start service
	if err := exec.Command("sc.exe", "start", serviceName).Run(); err != nil {
		return fmt.Errorf("failed to start service: %v", err)
	}

	console.Success("Magnesia Windows service installed and started successfully")
	return nil
}

func Uninstall() {
	serviceName := "MagnesiaAgent"
	_ = exec.Command("sc.exe", "stop", serviceName).Run()
	_ = exec.Command("sc.exe", "delete", serviceName).Run()
	_ = os.Remove(`C:\Program Files\Magnesia\magnesia.exe`)
	_ = os.RemoveAll(config.Dir())
	console.Warn("Magnesia Windows service removed")
}
