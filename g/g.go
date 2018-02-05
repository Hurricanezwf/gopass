package g

import (
	"os/user"
	"path/filepath"
	"runtime"

	"github.com/Hurricanezwf/gopass/log"
)

// 配置部分
var (
	// 配置路径
	ConfigPath string

	// 日志记录方式, "none", "console", "file"
	LogWay string

	// 日志级别
	LogLevel string

	// 日志文件路径
	LogFile string

	// 元数据存储目录
	MetaDir string
)

func init() {
	switch runtime.GOOS {
	case "linux":
		user, err := user.Current()
		if err != nil {
			panic("Get current user failed, " + err.Error())
		}
		if user == nil {
			panic("Get current user failed, user is nil")
		}
		LogWay = "none"
		MetaDir = filepath.Join(user.HomeDir, ".gopass")
	case "darwin":
		// TODO:
	case "windows":
		// TODO:
	}
}

func loadConf() {
	//ini, err := config.NewConfig("ini", "")
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
