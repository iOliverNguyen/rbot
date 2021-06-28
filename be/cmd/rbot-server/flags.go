package main

import (
	"flag"
	"fmt"
	"os"
)

const usage = "Usage: rbot-server [args ...]"

var flConfigFile = ""
var flFlowFile = ""
var flStateFile = ""
var flHelp = false

func initFlags() {
	flag.StringVar(&flConfigFile, "config-file", "", "path to config file")
	flag.StringVar(&flFlowFile, "flow-file", "./rbot-flow-data.json", "path to flow data file")
	flag.StringVar(&flStateFile, "state-file", "./rbot-state-data.json", "path to state data file")
	flag.BoolVar(&flHelp, "help", false, "")
	flag.Parse()

	fl := flag.CommandLine
	if flHelp {
		_, _ = fmt.Fprintf(fl.Output(), usage)
		fl.PrintDefaults()
		os.Exit(2)
	}
}
