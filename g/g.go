package g

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"

	"github.com/Hurricanezwf/gopass/log"
	"github.com/astaxie/beego/config"
)

// 配置部分
var (
	// 日志记录方式, "none", "console", "file"
	LogWay string

	// 日志级别
	LogLevel string

	// 日志文件路径
	LogFile string

	// 元数据存储目录
	MetaDir string
)

func LoadConf(configfile string) error {
	conf, err := config.NewConfig("ini", configfile)
	if err != nil {
		return err
	}

	// Log相关
	if LogWay = conf.String("LogWay"); len(LogWay) <= 0 {
		// use default
		LogWay = "none"
	}
	switch LogWay {
	case LogWayFile:
		goto TAG1
	case LogWayConsole:
		goto TAG2
	case LogWayNone:
		goto TAG3
	default:
		return errors.New("Invalid LogWay")
	}

TAG1:
	if LogFile = conf.String("LogFile"); len(LogFile) <= 0 {
		return errors.New("Missing LogFile")
	}

TAG2:
	if LogLevel = conf.String("LogLevel"); len(LogLevel) <= 0 {
		return errors.New("Missing LogLevel")
	}
	switch LogLevel {
	case "error", "warn", "notice", "info", "debug":
	default:
		return errors.New("Invalid LogLevel")
	}

TAG3:
	initLog()

	// Core相关
	if MetaDir = conf.String("MetaDir"); len(MetaDir) <= 0 {
		// use default
		switch runtime.GOOS {
		case "linux", "darwin":
			user, err := user.Current()
			if err != nil {
				return fmt.Errorf("Get current user failed, %v", err)
			}
			if user == nil {
				return errors.New("Get current user failed, user is nil")
			}
			MetaDir = filepath.Join(user.HomeDir, ".gopass")
		//case "darwin":
		// TODO:
		case "windows":
			// TODO:
		}

	}
	info, err := os.Stat(MetaDir)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("MetaDir(%s) doesnt' existed, please create it by yourself", MetaDir)
		}
		return err
	}
	if info.IsDir() == false {
		return fmt.Errorf("MetaDir(%s) is not a dir", MetaDir)
	}

	log.Info("========================= GOPASS CONFIG ==========================")
	log.Info("% -24s : %s", LogWay)
	log.Info("% -24s : %s", LogLevel)
	log.Info("% -24s : %s", LogFile)
	log.Info("% -24s : %s", MetaDir)
	log.Info("==================================================================")

	return nil
}

func initLog() {
	if LogWay != LogWayNone {
		logConf := &log.LogConf{}
		logConf.SetLogLevel(LogLevel)
		logConf.SetLogWay(LogWay)
		logConf.SetLogFile(LogFile)
		log.Init(logConf)
	}
}
