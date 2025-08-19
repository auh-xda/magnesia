//go:build windows
// +build windows

package services

import (
	"fmt"

	"github.com/auh-xda/magnesia/helpers/console"
	"golang.org/x/sys/windows/svc/mgr"
)

func ListServices() ([]WindowsService, error) {

	var services []WindowsService

	console.Info("Getting services for windows")

	// Connect to the Windows service manager
	m, err := mgr.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to service manager: %w", err)
	}
	defer m.Disconnect()

	// List all service names
	serviceNames, err := m.ListServices()
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}

	for _, name := range serviceNames {
		s, err := m.OpenService(name)
		if err != nil {
			continue // skip services we can't open
		}
		defer s.Close()

		// Query config and status
		config, _ := s.Config()
		status, _ := s.Query()

		state := "stopped"
		if status.State == 4 { // SERVICE_RUNNING
			state = "running"
		}

		services = append(services, WindowsService{
			Name:        name,
			DisplayName: config.DisplayName,
			Status:      state,
			StartType:   startTypeToString(config.StartType),
		})
	}

	return services, nil
}

// helper to convert start type to readable string
func startTypeToString(t uint32) string {
	switch t {
	case 2:
		return "automatic"
	case 3:
		return "manual"
	case 4:
		return "disabled"
	default:
		return "unknown"
	}
}
