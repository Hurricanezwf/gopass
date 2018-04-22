package app

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Hurricanezwf/gopass/crypt"
	"github.com/Hurricanezwf/gopass/g"
	"github.com/Hurricanezwf/gopass/log"
	"github.com/Hurricanezwf/gopass/meta"
	"github.com/Hurricanezwf/gopass/password"
	"github.com/Hurricanezwf/gopass/ui"
	term "github.com/howeyc/gopass"
	cli "gopkg.in/urfave/cli.v2"
)

const VERSION = "0.1.0"

func Run() {
	app := &cli.App{
		Version:     VERSION,
		Name:        "gopass",
		Usage:       "A tool for managing your password in terminal",
		UsageText:   "gopass [command] [-c ConfigFile]",
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
				UsageText: "gopass add [-c ConfigFile] [key] [password] ",
				Action:    AddAction,
				Flags: []cli.Flag{
					&cli.PathFlag{Name: "c", Usage: "Specify config file's path", Value: ""},
				},
			},
			// del password
			&cli.Command{
				Name:      g.CMDDel,
				Usage:     "delete password from gopass",
				UsageText: "gopass del [-c ConfigFile] [key]",
				Action:    DelAction,
				Flags: []cli.Flag{
					&cli.PathFlag{Name: "c", Usage: "Specify config file's path", Value: ""},
				},
			},
			// update password
			&cli.Command{
				Name:      g.CMDUpdate,
				Usage:     "update password into gopass",
				UsageText: "gopass update [-c ConfigFile] [key] [NewPassword] ",
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
			// change SecretKey
			&cli.Command{
				Name:      g.CMDChSK,
				Usage:     "change auth SecretKey which you provide for authentication when init the app",
				UsageText: "gopass chsk [-c ConfigFile]",
				Action:    ChSKAction,
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
		msg := fmt.Sprintf("\033[31m%s\033[0m", err.Error())
		fmt.Println(msg)
		log.Error(msg)
	}
}

func AddAction(c *cli.Context) error {
	var (
		err  error
		args []string
		sk   []byte
	)

	if args = c.Args().Slice(); len(args) != 2 {
		cli.ShowCommandHelp(c, g.CMDAdd)
		return nil
	}

	if err = actionInit(c); err != nil {
		return err
	}

	if sk, err = auth(); err != nil {
		return fmt.Errorf("Auth failed, %v", err)
	}

	if err = password.Add(sk, []byte(args[0]), []byte(args[1])); err != nil {
		return fmt.Errorf("Add password failed, %v", err)
	}
	fmt.Printf("\033[32mAdd OK\033[0m\n")
	return nil
}

func DelAction(c *cli.Context) error {
	var (
		err  error
		args []string
		sk   []byte
	)

	if args = c.Args().Slice(); len(args) != 1 {
		cli.ShowCommandHelp(c, g.CMDDel)
		return nil
	}

	if err = actionInit(c); err != nil {
		return err
	}

	if sk, err = auth(); err != nil {
		return fmt.Errorf("Auth failed, %v", err)
	}

	if err = password.Del(sk, []byte(args[0])); err != nil {
		return fmt.Errorf("Del password for key(%s) failed, %v", args[0], err)
	}
	fmt.Printf("\033[32mDel OK\033[0m\n")
	return nil
}

func UpdateAction(c *cli.Context) error {
	var (
		err  error
		args []string
		sk   []byte
	)

	if args = c.Args().Slice(); len(args) != 2 {
		cli.ShowCommandHelp(c, g.CMDUpdate)
		return nil
	}

	if err = actionInit(c); err != nil {
		return err
	}

	if sk, err = auth(); err != nil {
		return fmt.Errorf("Auth failed, %v", err)
	}

	if err = password.Update(sk, []byte(args[0]), []byte(args[1])); err != nil {
		return fmt.Errorf("Update password for key(%s) failed, %v", args[0], err)
	}
	fmt.Printf("\033[32mUpdate OK\033[0m\n")
	return nil
}

func GetAction(c *cli.Context) error {
	var (
		err error
		sk  []byte
	)

	if err = actionInit(c); err != nil {
		return err
	}

	if sk, err = auth(); err != nil {
		return fmt.Errorf("Auth failed, %v", err)
	}

	if err = ui.Open(sk); err != nil {
		return fmt.Errorf("Open UI failed, %v", err)
	}
	defer ui.Close()

	return nil
}

func ChSKAction(c *cli.Context) error {
	var (
		err   error
		sk    []byte
		newSK []byte
	)

	if err = actionInit(c); err != nil {
		return err
	}

	if sk, err = auth("Old SecretKey: "); err != nil {
		return fmt.Errorf("Auth failed, %v", err)
	}

	if newSK, err = inputSKTwice("New SecretKey: ", "Confirm      : "); err != nil {
		return err
	}

	if err = password.ChangeSK(sk, newSK); err != nil {
		return fmt.Errorf("Change failed, Err: %v. Please refer log for more details.", err)
	}
	fmt.Printf("\033[32mChange OK\033[0m\n")
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

// 返回用户输入的SK
func auth(prompt ...string) ([]byte, error) {
	var (
		err          error
		skReserved   []byte
		skEntered    []byte
		skEnteredEnc []byte
	)

	skReserved, err = password.GetAuthSK()
	if err != nil {
		if err != meta.ErrNotExist {
			return nil, fmt.Errorf("Get auth sk failed, %v", err)
		}
		// init the app
		fmt.Printf("Please init the app when you login firstly.\n")
		fmt.Printf("Your reserved SecretKey will be used to encrypt your content and provide authentications.\n\n")
		skEntered, err = inputSKTwice("SecretKey: ", "Confirm  : ")
		if err != nil {
			return nil, fmt.Errorf("Init auth sk failed, %v", err)
		}

		// save to db
		if skReserved, err = password.InitAuthSK(skEntered); err != nil {
			return nil, fmt.Errorf("Init auth sk failed, %v", err)
		}
		skEnteredEnc = skReserved // skReserved is encrypted
	}

	// enter
	if len(skEntered) <= 0 {
		p := "SecretKey: "
		if len(prompt) > 0 {
			p = prompt[0]
		}
		if skEntered, err = askForSK(p); err != nil {
			return nil, fmt.Errorf("Ask for sk failed, %v", err)
		}
		skEnteredEnc = crypt.EncryptSK(skEntered)
	}

	// compare
	if bytes.Compare(skEnteredEnc, skReserved) != 0 {
		return nil, fmt.Errorf("SecretKey not match")
	}

	return skEntered, nil
}

func inputSKTwice(firstPrompt, againPrompt string) ([]byte, error) {
	fmt.Printf(firstPrompt)
	sk1, err := term.GetPasswdMasked()
	if err != nil {
		return nil, err
	}
	if len(sk1) <= 0 {
		return nil, fmt.Errorf("Empty SecretKey input")
	}

	fmt.Printf(againPrompt)
	sk2, err := term.GetPasswdMasked()
	if err != nil {
		return nil, err
	}

	if bytes.Compare(sk1, sk2) != 0 {
		return nil, fmt.Errorf("SecretKey didn't equal")
	}

	return sk1, nil
}

func askForSK(prompt string) ([]byte, error) {
	fmt.Printf(prompt)
	return term.GetPasswdMasked()
}
