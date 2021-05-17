// Copyright (c) 2021, Athonet S.r.l. All rights reserved.
// riccardo.manfrin@athonet.com

package main

import (
	"flag"

	"github.com/riccardomanfrin/netconfd/logger"
)

var configfile = flag.String("config", "netconfd.json", "Path to netconfd configuration file")
var logfile = flag.String("log", "", "Path to netconfd log file (default to syslog)")
var skipbootconfig = flag.Bool("skipbootconfig", false, "Skip initial startup config patching")

func main() {

	flag.Parse()

	logger.LoggerInit(*logfile)
	logger.LoggerSetLevel("INF")

	mgr := NewManager()
	mgr.LoadConfig(configfile)
	mgr.Start()

	for {
		select {}
	}
}
