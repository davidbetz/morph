package main

import (
	"flag"
	"os"
	"strconv"

	"github.com/davidbetz/morph/internal/parser"
	"github.com/davidbetz/morph/internal/targets"
	"github.com/davidbetz/morph/internal/util"
)

var verbose bool

func main() {
	verbose, _ = strconv.ParseBool(os.Getenv("VERBOSE"))
	platformPtr := flag.String("target", "", "text|jsonl|aws|gcp|azure|mssql")
	modePtr := flag.String("mode", "", "gnt|wlc")
	stylePtr := flag.String("style", "", "english|hebrew")
	flag.Parse()
	mode := *modePtr
	if len(mode) == 0 {
		util.Errorf("-mode is required: gnt|wlc")
	}
	platformString := *platformPtr
	if len(platformString) == 0 {
		util.Errorf("-targets is required: text|jsonl|aws|gcp|azure|mssql")
	}
	targets := targets.Create(platformString)
	err := targets.ValidateCloudConfig()
	if err != nil {
		util.Errorf(err.Error())
	}
	activeParser := parser.Create(mode, *stylePtr)
	err = activeParser.Count(targets)
	if err != nil {
		util.Errorf(err.Error())
	}
}
