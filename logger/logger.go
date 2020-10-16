package logger

import (
	"os"

	logging "github.com/op/go-logging"
)

// Log logging facility
var Log = logging.MustGetLogger("example")

var lfb logging.LeveledBackend
var dfltLogLevel = logging.WARNING

//LoggerSetLevel Configures the log level
func LoggerSetLevel(level string) int {
	switch level {
	case "DBG":
		lfb.SetLevel(logging.DEBUG, "")
	case "INF":
		lfb.SetLevel(logging.INFO, "")
	case "WRN":
		lfb.SetLevel(logging.WARNING, "")
	case "ERR":
		lfb.SetLevel(logging.ERROR, "")
	case "CRI":
		lfb.SetLevel(logging.CRITICAL, "")
	}

	return 0
}
func Fatal(err error) {
	Log.Warning(err.Error())
}

// LoggerInit initializes the log functionality (to be invoked on bootstrap)
func LoggerInit(file string) {
	var res string

	var format = logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{level:.4s} %{id:03x} [%{shortfile}] ▶ %{color:reset} %{message}`,
	)
	if file != "" {
		if file == "-" {
			res = "Logging to stdout"
			b := logging.NewLogBackend(os.Stdout, "", 0)
			f := logging.NewBackendFormatter(b, format)
			lfb = logging.AddModuleLevel(f)
			lfb.SetLevel(dfltLogLevel, "")
			logging.SetBackend(lfb)
		} else {
			logfile, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
			if err != nil {
				res = "Failed to open logfile [" + file + "]; defaulting to stdout"
				logfile = os.Stdout
			} else {
				res = "Logging to logfile [" + file + "]"
			}
			b := logging.NewLogBackend(logfile, "", 0)
			f := logging.NewBackendFormatter(b, format)
			lfb = logging.AddModuleLevel(f)
			lfb.SetLevel(dfltLogLevel, "")
			logging.SetBackend(lfb)
		}
	} else {
		var b, err = logging.NewSyslogBackend("netconfd")
		if err != nil {
			res = "Failed to open channel towards syslog, defaulting to stdout"
			var b = logging.NewLogBackend(os.Stdout, "", 0)
			f := logging.NewBackendFormatter(b, format)
			lfb = logging.AddModuleLevel(f)
			lfb.SetLevel(dfltLogLevel, "")
			logging.SetBackend(lfb)
		} else {
			f := logging.NewBackendFormatter(b, format)
			lfb = logging.AddModuleLevel(f)
			lfb.SetLevel(dfltLogLevel, "")
			logging.SetBackend(lfb)
			res = "Logging to syslog"
		}

	}

	Log.Info(res)
}
