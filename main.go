package main

import (
	"flag"

	"github.com/Hurricanezwf/gopass/log"
	"github.com/Hurricanezwf/gopass/ui"
)

var (
	logFile string
)

func init() {
	flag.StringVar(&logFile, "l", "./gopass", "-l logfile")
	flag.Parse()
}

func main() {
	initLog()

	ui.Open()
}

func initLog() {
	logConf := log.DefaultLogConf()
	logConf.SetLogLevel("warn")
	logConf.SetLogWay(log.LogWayFile)
	logConf.SetLogFile(logFile)
	log.Init(logConf)
}
