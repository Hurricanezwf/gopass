package app

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Hurricanezwf/gopass/g"
	"github.com/Hurricanezwf/gopass/password"
	"github.com/Hurricanezwf/gopass/ui"
	cli "gopkg.in/urfave/cli.v2"
)

const VERSION = "0.0.4"

func Run() {
	app := &cli.App{
		Version:     VERSION,
		Name:        "gopass",
		Usage:       "A tool for managing your password in terminal",
		UsageText:   "gopass [command]",
		HideVersion: true,
		HideHelp:    true,
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "h", Usage: "Show help", Hidden: true},
		},
		Commands: []*cli.Command{
			// add password
			&cli.Command{
				Name:      g.CMDAdd,
				Usage:     "add password into gopass",
				UsageText: "gopass add [key] [password] [-c ConfigFile]",
				Action:    AddAction,
				Flags: []cli.Flag{
					&cli.PathFlag{Name: "c", Usage: "Specify config file's path", Value: ""},
				},
			},
			// del password
			&cli.Command{
				Name:      g.CMDDel,
				Usage:     "delete password from gopass",
				UsageText: "gopass del [key] [-c ConfigFile]",
				Action:    DelAction,
				Flags: []cli.Flag{
					&cli.PathFlag{Name: "c", Usage: "Specify config file's path", Value: ""},
				},
			},
			// update password
			&cli.Command{
				Name:      g.CMDUpdate,
				Usage:     "update password into gopass",
				UsageText: "gopass update [key] [new_password] [-c ConfigFile]",
				Action:    UpdateAction,
				Flags: []cli.Flag{
					&cli.PathFlag{Name: "c", Usage: "Specify config file's path", Value: ""},
				},
			},
			// get password
			&cli.Command{
				Name:      g.CMDGet,
				Usage:     "display ui to search and copy password",
				UsageText: "gopass ui [-c ConfigFile]",
				Action:    GetAction,
				Flags: []cli.Flag{
					&cli.PathFlag{Name: "c", Usage: "Specify config file's path", Value: ""},
				},
			},
			// show help
			&cli.Command{
				Name:   g.CMDHelp,
				Usage:  "show help",
				Action: HelpAction,
			},
			// show version
			&cli.Command{
				Name:   g.CMDVersion,
				Usage:  "show version",
				Action: VersionAction,
			},
		},
	}

	app.Before = func(c *cli.Context) error {
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

func AddAction(c *cli.Context) error {
	var (
		err  error
		args []string
	)

	if args = c.Args().Slice(); len(args) != 2 {
		cli.ShowCommandHelp(c, g.CMDAdd)
		return nil
	}

	if err = actionInit(c); err != nil {
		return err
	}

	if err = password.Add([]byte(args[0]), []byte(args[1])); err != nil {
		return fmt.Errorf("Add password failed, %v", err)
	}
	fmt.Printf("\033[32mAdd OK\033[0m\n")
	return nil
}

func DelAction(c *cli.Context) error {
	var (
		err  error
		args []string
	)

	if args = c.Args().Slice(); len(args) != 1 {
		cli.ShowCommandHelp(c, g.CMDDel)
		return nil
	}

	if err = actionInit(c); err != nil {
		return err
	}

	if err = password.Del([]byte(args[0])); err != nil {
		return fmt.Errorf("Del password for key(%s) failed, %v", args[0], err)
	}
	fmt.Printf("\033[32mDel OK\033[0m\n")
	return nil
}

func UpdateAction(c *cli.Context) error {
	var (
		err  error
		args []string
	)

	if args = c.Args().Slice(); len(args) != 2 {
		cli.ShowCommandHelp(c, g.CMDUpdate)
		return nil
	}

	if err = actionInit(c); err != nil {
		return err
	}

	if err = password.Update([]byte(args[0]), []byte(args[1])); err != nil {
		return fmt.Errorf("Update password for key(%s) failed, %v", args[0], err)
	}
	fmt.Printf("\033[32mUpdate OK\033[0m\n")
	return nil
}

func GetAction(c *cli.Context) error {
	var err error

	if err = actionInit(c); err != nil {
		return err
	}

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

func VersionAction(c *cli.Context) error {
	cli.ShowVersion(c)
	return nil
}

func actionInit(c *cli.Context) error {
	configFile := c.Path("c")
	if len(configFile) <= 0 {
		// 默认查找与二进制文件同目录的conf.ini
		path, err := exec.LookPath(os.Args[0])
		if err != nil {
			return fmt.Errorf("Find binary path failed, %v", err)
		}
		path, err = filepath.Abs(filepath.Dir(path))
		if err != nil {
			return fmt.Errorf("Convert path to abs failed, %v", err)
		}
		configFile = filepath.Join(path, "conf.ini")
	}

	if err := g.LoadConf(configFile); err != nil {
		return fmt.Errorf("Load config failed, %v. ConfigFile: %s", err, configFile)
	}
	return nil
}
