package config

var verbose bool

func SetVerbose(v bool) {
	verbose = v
}

func GetVerbose() bool {
	return verbose
}
