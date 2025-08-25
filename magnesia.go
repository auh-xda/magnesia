package main

import (
	"github.com/auh-xda/magnesia/config"
	"github.com/auh-xda/magnesia/console"
)

func (magnesia Magnesia) Installed() bool {
	_, err := config.ParseConfig()

	return nil == err
}

func (magnesia Magnesia) Info() {
	config, _ := config.ParseConfig()

	console.Table(config)
}
