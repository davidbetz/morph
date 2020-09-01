package main

import (
	"fmt"
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
		fmt.Println("MODE is required: gnt|wlc")
		return
	}
	err := validateCloudConfig()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	var parser parser
	if mode == "wlc" {
		parser = CreateWlc()
	} else {
		parser = CreateGnt()
	}
	err = parser.Process()
	if err != nil {
		fmt.Println(err.Error())
	}
}
