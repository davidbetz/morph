package main

import (
	"flag"
	"os"
	"strconv"

	"github.com/davidbetz/morph/internal/macula"
	"github.com/davidbetz/morph/internal/models"
	"github.com/davidbetz/morph/internal/parser"
	"github.com/davidbetz/morph/internal/targets"
	"github.com/davidbetz/morph/internal/util"
)

var verbose bool

func main() {
	verbose, _ = strconv.ParseBool(os.Getenv("VERBOSE"))
	platformPtr := flag.String("target", "", "text|jsonl|aws|gcp|azure|mssql")
	modePtr := flag.String("mode", "", "gnt|wlc")
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
	var activeParser models.Parser
	if mode == "gnt" {
		activeParser = parser.CreateGnt()
	} else if mode == "wlc" {
		activeParser = macula.CreateHebrewParser()
	}
	err = activeParser.Count(targets)
	if err != nil {
		util.Errorf(err.Error())
	}
}
