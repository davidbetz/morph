package main

import (
	"flag"
	"os"
	"strconv"

	"github.com/davidbetz/morph/internal/parser"
	"github.com/davidbetz/morph/internal/platform"
	"github.com/davidbetz/morph/internal/util"
)

var verbose bool

type activeParser interface {
	Process() error
}

func main() {
	verbose, _ = strconv.ParseBool(os.Getenv("VERBOSE"))
	modePtr := flag.String("mode", "", "gnt|wlc")
	stylePtr := flag.String("style", "", "english|hebrew")
	flag.Parse()
	mode := *modePtr
	if len(mode) == 0 {
		util.Errorf("-mode is required: gnt|wlc")
	}
	err := platform.ValidateCloudConfig()
	if err != nil {
		util.Errorf(err.Error())
	}
	var activeParser activeParser
	if mode == "wlc" {
		style := *stylePtr
		activeParser = parser.CreateWlc(style)
	} else {
		activeParser = parser.CreateGnt()
	}
	err = activeParser.Process()
	if err != nil {
		util.Errorf(err.Error())
	}
}
