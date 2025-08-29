//go:build windows
// +build windows

package installer

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/auh-xda/magnesia/config"
	"github.com/auh-xda/magnesia/console"
)

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open source: %w", err)
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("create target: %w", err)
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return fmt.Errorf("copy contents: %w", err)
	}

	return nil
}

func CreateService() error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %v", err)
	}

	targetBin := `C:\Program Files\Magnesia\magnesia.exe`

	// Ensure target directory exists
	if err := os.MkdirAll(filepath.Dir(targetBin), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Copy binary using Go-native function
	if err := copyFile(exePath, targetBin); err != nil {
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
