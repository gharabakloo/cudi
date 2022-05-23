package config

import "flag"

var V bool
var File string

func InitializeArg() {
	flag.BoolVar(&V, "v", false, "verbosity")
	flag.StringVar(&File, "config", "config.json", "config file path")
}
