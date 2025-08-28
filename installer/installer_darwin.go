//go:build darwin
// +build darwin

package installer

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/auh-xda/magnesia/config"
	"github.com/auh-xda/magnesia/console"
)

func CreateService() error {
	plistFile := `/Library/LaunchDaemons/com.magnesia.agent.plist`
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %v", err)
	}

	// Copy binary to /usr/local/bin
	targetBin := "/usr/local/bin/magnesia"
	if err := exec.Command("cp", exePath, targetBin).Run(); err != nil {
		return fmt.Errorf("failed to copy binary: %v", err)
	}
	if err := os.Chmod(targetBin, 0755); err != nil {
		return fmt.Errorf("chmod failed: %v", err)
	}

	// Create plist file
	plistContent := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple Computer//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>Label</key>
	<string>com.magnesia.agent</string>
	<key>ProgramArguments</key>
	<array>
		<string>%s</string>
	</array>
	<key>RunAtLoad</key>
	<true/>
	<key>KeepAlive</key>
	<true/>
	<key>StandardOutPath</key>
	<string>/var/log/magnesia.log</string>
	<key>StandardErrorPath</key>
	<string>/var/log/magnesia.log</string>
</dict>
</plist>`, targetBin)

	if err := os.WriteFile(plistFile, []byte(plistContent), 0644); err != nil {
		return fmt.Errorf("failed to write plist: %v", err)
	}

	// Load service
	if err := exec.Command("launchctl", "load", plistFile).Run(); err != nil {
		return fmt.Errorf("failed to load service: %v", err)
	}

	console.Success("Magnesia agent installed as launchd service (logs: /var/log/magnesia.log)")
	return nil
}

func Uninstall() {
	plistFile := `/Library/LaunchDaemons/com.magnesia.agent.plist`
	_ = exec.Command("launchctl", "unload", plistFile).Run()
	_ = os.Remove(plistFile)
	_ = os.Remove("/usr/local/bin/magnesia")
	_ = os.Remove("/var/log/magnesia.log")
	_ = os.RemoveAll(config.Dir())
	console.Warn("Magnesia service removed from macOS")
}
