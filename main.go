package main

import (
	"os"
	"strconv"
)

var verbose bool
var sourceFileLocation string

type parser interface {
	Process() error
}

func main() {
	sourceFileLocation = os.Getenv("SOURCE")
	verbose, _ = strconv.ParseBool(os.Getenv("VERBOSE"))
	mode := os.Getenv("MODE")
	if len(mode) == 0 {
		errorf("MODE is required: gnt|wlc")
	}
	err := validateCloudConfig()
	if err != nil {
		errorf(err.Error())
	}
	var parser parser
	if mode == "wlc" {
		parser = CreateWlc()
	} else {
		parser = CreateGnt()
	}
	err = parser.Process()
	if err != nil {
		errorf(err.Error())
	}
}
