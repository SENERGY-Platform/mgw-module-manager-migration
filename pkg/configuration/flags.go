package configuration

import "flag"

var ConfPath string

func ParseFlags() {
	flag.StringVar(&ConfPath, "config", "", "path to config JSON file")
	flag.Parse()
	return
}
