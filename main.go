package main

import (
	"flag"

	"gitlab.lan.athonet.com/core/netconfd/logger"
)

var configfile = flag.String("config", "netconfd.json", "Path to netconfd configuration file")
var logfile = flag.String("log", "", "Path to netconfd log file (default to syslog)")

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
