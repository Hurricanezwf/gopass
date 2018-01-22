package main

import (
	"fmt"
	"os"

	"github.com/Hurricanezwf/gopass/log"
	"github.com/Hurricanezwf/gopass/service"
	"github.com/Hurricanezwf/gopass/ui"
	cli "gopkg.in/urfave/cli.v2"
)

var (
	logFile string = "./gopass.log"
)

// Command list
const (
	ADD  = "add"
	GET  = "get"
	DEL  = "del"
	HELP = "help"
)

//func init() {
//	flag.StringVar(&logFile, "l", "./gopass.log", "-l logfile")
//	flag.Parse()
//	initLog()
//}

func main() {
	app := &cli.App{
		Version:     "0.0.1",
		Name:        "gopass",
		UsageText:   "gopass [add|del]",
		HideVersion: true,
		HideHelp:    true,
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "h", Usage: "show help", Hidden: true},
		},
		Commands: []*cli.Command{
			// add password
			&cli.Command{
				Name:      ADD,
				Usage:     "add password to manager",
				UsageText: "gopass add [key] [password]",
				Action:    AddAction,
			},
			// del password
			&cli.Command{
				Name:   DEL,
				Action: DelAction,
			},
			// get password
			&cli.Command{
				Name:   GET,
				Action: GetAction,
			},
			// show help
			&cli.Command{
				Name:   HELP,
				Action: HelpAction,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("\033[33m%v\033[0m\n", err)
	}
}

func AddAction(c *cli.Context) error {
	var (
		err  error
		args []string
	)

	if args = c.Args().Slice(); len(args) != 2 {
		cli.ShowCommandHelp(c, ADD)
		return nil
	}

	if err = service.Open(); err != nil {
		return fmt.Errorf("Open service failed, %v", err)
	}
	defer service.Close()

	if err = service.Passwd.Add([]byte(args[0]), []byte(args[1])); err != nil {
		return fmt.Errorf("Add password failed, %v", err)
	}

	fmt.Printf("\033[32mAdd OK\033[0m\n")
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

func HelpAction(c *cli.Context) error {
	cli.ShowAppHelp(c)
	return nil
}

func initLog() {
	logConf := &log.LogConf{}
	logConf.SetLogLevel("warn")
	logConf.SetLogWay(log.LogWayFile)
	logConf.SetLogFile(logFile)
	log.Init(logConf)
}
