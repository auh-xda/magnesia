//go:build darwin
// +build darwin

package services

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/auh-xda/magnesia/helpers/console"
)

func ListServices() ([]DarwinService, error) {
	var services []DarwinService

	console.Info("Getting services for darwin/mac")

	cmd := exec.Command("launchctl", "list")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run launchctl: %w", err)
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))
	firstLine := true
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		if firstLine { // skip header
			firstLine = false
			continue
		}
		if len(line) < 3 {
			continue
		}

		// PID
		pid := 0
		if line[0] != "-" {
			pid, _ = strconv.Atoi(line[0])
		}

		// Status (0 means running ok, nonzero means error)
		status := "stopped"
		if pid > 0 {
			status = "running"
		}

		services = append(services, DarwinService{
			Label:  line[2],
			Status: status,
			PID:    pid,
		})
	}

	return services, nil
}
