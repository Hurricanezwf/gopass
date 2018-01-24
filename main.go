package main

import (
	"fmt"
	"os"

	"github.com/Hurricanezwf/gopass/log"
	"github.com/Hurricanezwf/gopass/service"
	"github.com/Hurricanezwf/gopass/ui"
	cli "gopkg.in/urfave/cli.v2"
)

// Command list
const (
	ADD    = "add"
	GET    = "get"
	DEL    = "del"
	UPDATE = "update"
	HELP   = "help"
)

func main() {
	app := &cli.App{
		Version:     "0.0.1",
		Name:        "gopass",
		UsageText:   "gopass [add|del]",
		HideVersion: true,
		HideHelp:    true,
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "h", Usage: "Show help", Hidden: true},
			&cli.PathFlag{Name: "logfile", Aliases: []string{"l"}, Usage: "Specify log file's path", Value: "./gopass.log"},
		},
		Commands: []*cli.Command{
			// add password
			&cli.Command{
				Name:      ADD,
				Usage:     "add password to manager",
				UsageText: "gopass add $(key) $(password)",
				Action:    AddAction,
			},
			// del password
			&cli.Command{
				Name:      DEL,
				Usage:     "delete password to manager",
				UsageText: "gopass del $(key)",
				Action:    DelAction,
			},
			// update password
			&cli.Command{
				Name:      UPDATE,
				Usage:     "update password to manager",
				UsageText: "gopass update $(key) $(new_password)",
				Action:    UpdateAction,
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

	app.Before = func(c *cli.Context) error {
		initLog(c.Path("logfile"))
		return nil
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

	/*
		if err = service.Open(); err != nil {
			return fmt.Errorf("Open service failed, %v", err)
		}
		defer service.Close()

		if err = service.Passwd.Add([]byte(args[0]), []byte(args[1])); err != nil {
			return fmt.Errorf("Add password failed, %v", err)
		}
	*/
	if err = service.AddPassword([]byte(args[0]), []byte(args[1])); err != nil {
		return fmt.Errorf("Add password failed, %v", err)
	}
	service.CloseAll()
	fmt.Printf("\033[32mAdd OK\033[0m\n")
	return nil
}

func DelAction(c *cli.Context) error {
	var (
		err  error
		args []string
	)

	if args = c.Args().Slice(); len(args) != 1 {
		cli.ShowCommandHelp(c, DEL)
		return nil
	}

	/*
		if err = service.Open(); err != nil {
			return fmt.Errorf("Open service failed, %v", err)
		}
		defer service.Close()

		if err = service.Passwd.Del([]byte(args[0])); err != nil {
			return fmt.Errorf("Del password for key(%s) failed, %v", args[0], err)
		}
	*/
	if err = service.DelPassword([]byte(args[0])); err != nil {
		return fmt.Errorf("Del password for key(%s) failed, %v", args[0], err)
	}
	service.CloseAll()
	fmt.Printf("\033[32mDel OK\033[0m\n")
	return nil
}

func UpdateAction(c *cli.Context) error {
	var (
		err  error
		args []string
	)

	if args = c.Args().Slice(); len(args) != 2 {
		cli.ShowCommandHelp(c, UPDATE)
		return nil
	}

	/*
		if err = service.Open(); err != nil {
			return fmt.Errorf("Open service failed, %v", err)
		}
		defer service.Close()

		if err = service.Passwd.Update([]byte(args[0]), []byte(args[1])); err != nil {
			return fmt.Errorf("Update password for key(%s) failed, %v", args[0], err)
		}
	*/

	if err = service.UpdatePassword([]byte(args[0]), []byte(args[1])); err != nil {
		return fmt.Errorf("Update password for key(%s) failed, %v", args[0], err)
	}
	service.CloseAll()
	fmt.Printf("\033[32mUpdate OK\033[0m\n")
	return nil
}

func GetAction(c *cli.Context) error {
	var err error

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

func initLog(logFile string) {
	if len(logFile) <= 0 {
		logFile = "./gopass.log"
	}

	logConf := &log.LogConf{}
	logConf.SetLogLevel("debug")
	logConf.SetLogWay(log.LogWayFile)
	logConf.SetLogFile(logFile)
	log.Init(logConf)
}
