package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Hurricanezwf/gopass/log"
	"github.com/Hurricanezwf/gopass/service"
	"github.com/Hurricanezwf/gopass/ui"
	cli "gopkg.in/urfave/cli.v2"
)

var (
	logFile string
)

// Command list
const (
	ADD = "add"
	GET = "get"
	DEL = "del"
)

func init() {
	flag.StringVar(&logFile, "l", "./gopass", "-l logfile")
	flag.Parse()
	initLog()
}

func main() {
	app := &cli.App{
		Version:   "0.0.1",
		Name:      "gopass",
		UsageText: "gopass [add|del]",
		Commands:  make([]*cli.Command, 3),
	}
	app.Commands[0] = &cli.Command{
		Name:      ADD,
		Usage:     "save password",
		UsageText: "gopass add [key] [password]",
		Action:    AddAction,
	}
	app.Commands[1] = &cli.Command{
		Name:   DEL,
		Action: DelAction,
	}
	app.Commands[2] = &cli.Command{
		Name:   GET,
		Action: GetAction,
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

func AddAction(c *cli.Context) error {
	fmt.Printf("Add action\n")

	args := c.Args().Slice()
	for i, a := range args {
		fmt.Printf("[%d] %s\n", i, a)
	}

	if c.Args().Len() != 2 {
		cli.ShowCommandHelp(c, ADD)
	}
	return nil
}

func DelAction(c *cli.Context) error {
	fmt.Printf("Del action\n")
	return nil
}

func GetAction(c *cli.Context) error {
	var err error
	if err = service.Open(); err != nil {
		return fmt.Errorf("Open service failed, %v", err)
	}
	defer service.Close()

	if err = ui.Open(); err != nil {
		return fmt.Errorf("Open UI failed, %v", err)
	}
	defer ui.Close()

	return nil
}

func initLog() {
	logConf := &log.LogConf{}
	logConf.SetLogLevel("warn")
	logConf.SetLogWay(log.LogWayFile)
	logConf.SetLogFile(logFile)
	log.Init(logConf)
}
