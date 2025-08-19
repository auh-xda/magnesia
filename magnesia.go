package main

import "github.com/auh-xda/magnesia/helpers/console"

func (Magnesia) Installed() bool {
	_, err := (Config{}).Parse()

	return nil == err
}

func (Magnesia) Info() {
	config, _ := (Config{}).Parse()

	console.Table(config)
}
