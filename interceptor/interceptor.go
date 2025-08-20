package interceptor

import (
	"fmt"

	"github.com/auh-xda/magnesia/console"
)

func GetServices() {
	services, err := ListServices()

	if err != nil {
		console.Error(err.Error())
		return
	}

	console.Success(fmt.Sprintf("%d service fetched successfully", len(services)))
}

func BatteryInfo() PowerInfo {

	powerInfo, _ := GetPowerInfo()

	return powerInfo
}

func GetCpuDetails() CPUInfo {

	info, err := GetCPUInfo()

	if err != nil {
		console.Error(err.Error())
	}

	return info
}
