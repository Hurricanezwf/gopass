package g

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

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
	LogWay = conf.String("LogWay")
	LogWay = strings.ToLower(LogWay)
	switch LogWay {
	case LogWayFile:
		{
			if LogFile = conf.String("LogFile"); len(LogFile) <= 0 {
				return errors.New("Missing `LogFile` field")
			}

			LogLevel = conf.String("LogLevel")
			LogLevel = strings.ToLower(LogLevel)
			switch LogLevel {
			case "error", "warn", "notice", "info", "debug":
			default:
				LogLevel = "warn"
			}

			initLog()
		}
	case LogWayDisabled:
		// do nothing
	default:
		return errors.New("Invalid LogWay")
	}

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
		case "windows":
			// TODO:
		}

	}
	//info, err := os.Stat(MetaDir)
	//if err != nil {
	//	if os.IsNotExist(err) {
	//		return fmt.Errorf("MetaDir(%s) doesnt' existed, please create it by yourself", MetaDir)
	//	}
	//	return err
	//}
	//if info.IsDir() == false {
	//	return fmt.Errorf("MetaDir(%s) is not a dir", MetaDir)
	//}

	log.Info("========================= GOPASS CONFIG ==========================")
	log.Info("% -24s : %s", "LogWay", LogWay)
	log.Info("% -24s : %s", "LogLevel", LogLevel)
	log.Info("% -24s : %s", "LogFile", LogFile)
	log.Info("% -24s : %s", "MetaDir", MetaDir)
	log.Info("==================================================================")

	return nil
}

func initLog() {
	if LogWay == LogWayDisabled {
		return
	}

	logDir := filepath.Dir(LogFile)
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		err = fmt.Errorf("Create log dir %s failed, %v", logDir, err)
		panic(err)
	}

	logConf := &log.LogConf{}
	logConf.SetLogLevel(LogLevel)
	logConf.SetLogWay(LogWay)
	logConf.SetLogFile(LogFile)
	log.Init(logConf)
}
