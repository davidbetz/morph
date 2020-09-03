package main

import (
	"flag"
	"os"
	"strconv"
)

var verbose bool

type parser interface {
	Process() error
}

func main() {
	verbose, _ = strconv.ParseBool(os.Getenv("VERBOSE"))
	modePtr := flag.String("mode", "", "gnt|wlc")
	stylePtr := flag.String("style", "", "english|hebrew")
	flag.Parse()
	mode := *modePtr
	if len(mode) == 0 {
		errorf("-mode is required: gnt|wlc")
	}
	err := validateCloudConfig()
	if err != nil {
		errorf(err.Error())
	}
	var parser parser
	if mode == "wlc" {
		style := *stylePtr
		parser = CreateWlc(style)
	} else {
		parser = CreateGnt()
	}
	err = parser.Process()
	if err != nil {
		errorf(err.Error())
	}
}
