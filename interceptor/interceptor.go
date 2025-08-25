package interceptor

import (
	"fmt"

	"github.com/auh-xda/magnesia/console"
	"github.com/auh-xda/magnesia/nats"
)

func GetServices() {
	services, err := ListServices()

	if err != nil {
		console.Error(err.Error())
		return
	}

	nats.SendData(services, "services")

	console.Success(fmt.Sprintf("%d service fetched successfully", len(services)))
}

func BatteryInfo(sendToNats bool) PowerInfo {

	powerInfo, _ := GetPowerInfo()

	if sendToNats {
		nats.SendData(powerInfo, "power_info")
	}

	return powerInfo
}

func GetCpuDetails() CPUInfo {

	info, err := GetCPUInfo()

	if err != nil {
		console.Error(err.Error())
	}

	return info
}

func InstalledSoftwareList() {

	sw, err := Installations()

	if err != nil {
		console.Error("failed to query installed software list")
	}

	nats.SendData(sw, "installations")

	console.Success(fmt.Sprintf("%d softwares are there ", len(sw)))
}
