package log

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/astaxie/beego/logs"
)

var (
	LogWayConsole = logs.AdapterConsole
	LogWayFile    = logs.AdapterFile
)

////////////////////////////////////////////////////////////////////////////
type LogConf struct {
	logLevel int
	logWay   string
	logFile  string
}

func DefaultLogConf() *LogConf {
	return &LogConf{
		logLevel: logs.LevelInfo,
		logWay:   logs.AdapterConsole,
		logFile:  "",
	}
}

func (lc *LogConf) SetLogLevel(level string) *LogConf {
	var lv int
	switch level {
	case "error":
		lv = logs.LevelError
	case "warn":
		lv = logs.LevelWarning
	case "notice":
		lv = logs.LevelNotice
	case "info":
		lv = logs.LevelInformational
	case "debug":
		lv = logs.LevelDebug
	default:
		lv = logs.LevelInfo
	}

	lc.logLevel = lv
	return lc
}

func (lc *LogConf) LogLevel() string {
	switch lc.logLevel {
	case logs.LevelDebug:
		return "debug"
	case logs.LevelInformational:
		return "info"
	case logs.LevelNotice:
		return "notice"
	case logs.LevelWarning:
		return "warn"
	case logs.LevelError:
		return "error"
	}
	return ""
}

func (lc *LogConf) SetLogWay(way string) *LogConf {
	if way == logs.AdapterConsole || way == logs.AdapterFile {
		lc.logWay = way
	}
	return lc
}

func (lc *LogConf) SetLogFile(file string) *LogConf {
	absPath, err := filepath.Abs(file)
	if err != nil {
		log.Printf("Invalid filepath, %v\n", err)
		return lc
	}
	lc.logFile = absPath
	return lc
}

///////////////////////////////////////////////////////////////////////////////
type Logger struct {
	log *logs.BeeLogger

	cfg *LogConf
}

func NewLogger() *Logger {
	logger := logs.NewLogger(1000)
	logger.EnableFuncCallDepth(true)
	logger.SetLogFuncCallDepth(4)

	return &Logger{
		log: logger,
		cfg: DefaultLogConf(),
	}
}

// Apply applies all config to BeeLogger
func (l *Logger) Apply(lc *LogConf) *Logger {
	if lc != nil {
		l.cfg = lc
		config := fmt.Sprintf("{\"filename\":\"%s\",\"level\":%d}", lc.logFile, lc.logLevel)
		l.log.DelLogger(lc.logWay)
		l.log.SetLogger(lc.logWay, config)
	}
	return l
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.log.Error(format, v...)
}

func (l *Logger) Warn(format string, v ...interface{}) {
	l.log.Warn(format, v...)
}

func (l *Logger) Notice(format string, v ...interface{}) {
	l.log.Notice(format, v...)
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.log.Info(format, v...)
}

func (l *Logger) Debug(format string, v ...interface{}) {
	l.log.Debug(format, v...)
}
