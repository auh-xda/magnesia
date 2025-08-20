//go:build linux
// +build linux

package interceptor

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/auh-xda/magnesia/console"
)

func ListServices() ([]LinuxService, error) {
	var services []LinuxService

	console.Info("Getting services for Linux")

	// Run systemctl to get list of services
	cmd := exec.Command("systemctl",
		"list-units",
		"--type=service",
		"--all",
		"--no-legend",
		"--no-pager",
		"--plain")

	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run systemctl: %w", err)
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		if len(line) < 4 {
			continue
		}

		unit := line[0] // e.g. ssh.service

		active := line[2] // "active"/"inactive"/"failed"
		description := strings.Join(line[4:], " ")

		state := "stopped"
		if active == "active" {
			state = "running"
		}

		services = append(services, LinuxService{
			Name:        strings.TrimSuffix(unit, ".service"),
			Description: description,
			Status:      state,
		})
	}

	return services, nil
}

func GetInfo() PowerInfo {
	return PowerInfo{
		Vendor:   readFile("/sys/class/power_supply/BAT0/manufacturer"),
		Model:    readFile("/sys/class/power_supply/BAT0/model_name"),
		Serial:   readFile("/sys/class/power_supply/BAT0/serial_number"),
		Status:   readFile("/sys/class/power_supply/BAT0/status"),
		Capacity: readFile("/sys/class/power_supply/BAT0/capacity"),
	}
}

func readFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return "N/A"
	}
	return strings.TrimSpace(string(data))
}
