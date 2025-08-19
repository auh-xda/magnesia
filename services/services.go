package services

import (
	"fmt"

	"github.com/auh-xda/magnesia/helpers/console"
)

func GetServiceList() any {

	services, err := ListServices()

	if err != nil {
		console.Log(err)
	}

	console.Success(fmt.Sprintf("%d services running", len(services)))

	return services
}
