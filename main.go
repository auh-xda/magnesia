//       *******************************************************************
//       *     __  __                                        _             *
//       *    |  \/  |   __ _    __ _   _ __     ___   ___  (_)   __ _     *
//       *    | |\/| |  / _` |  / _` | | '_ \   / _ \ / __| | |  / _` |    *
//       *    | |  | | | (_| | | (_| | | | | | |  __/ \__ \ | | | (_| |    *
//       *    |_|  |_|  \__,_|  \__, | |_| |_|  \___| |___/ |_|  \__,_|    *
//       *                      |___/                                      *
//       *******************************************************************

package main

import (
	"flag"
	"fmt"

	"github.com/auh-xda/magnesia/console"
	"github.com/auh-xda/magnesia/interceptor"
)

const (
	authEndpoint   = "/c/0019-b33e-48e6-ba83"
	natsWsEndpoint = "nats://127.0.0.1:4222"
	version        = "0.1.0"
)

func main() {
	action := flag.String("action", "install", "Action to perform: install, update, or remove the Magnesia agent")
	auth_token := flag.String("auth_token", "", "Authentication token provided by the server")
	api_key := flag.String("api_key", "", "API key for your account")
	client_id := flag.String("client_id", "", "Unique client identifier")
	client_secret := flag.String("client_secret", "", "Client secret used for secure authentication")

	flag.Parse()

	magnesia := Magnesia{
		AuthToken:    *auth_token,
		ApiKey:       *api_key,
		ClientID:     *client_id,
		ClientSecret: *client_secret,
	}

	if *action != "install" && !magnesia.Installed() {
		console.Error("Magnesia not installed")
		return
	}

	switch *action {
	case "install":
		magnesia.Install()

	case "intercept":
		magnesia.Intercept()

	case "plist":
		magnesia.ProcessList()

	case "services":
		interceptor.GetServices()

	case "power":
		p := interceptor.BatteryInfo()
		console.Log(p)

	case "info":
		magnesia.Info()

	default:
		console.Error(fmt.Sprintf("Magnesia is not aware of this action (i.e %s)", *action))
	}
}
