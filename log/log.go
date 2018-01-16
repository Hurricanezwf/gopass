package log

var (
	l *Logger
)

func init() {
	conf := DefaultLogConf()
	conf.SetLogLevel("debug")
	conf.SetLogWay(LogWayFile)
	conf.SetLogFile("./gopass.log")

	l = NewLogger().Apply(conf)
}

func Error(format string, v ...interface{}) {
	l.Error(format, v...)
}

func Warn(format string, v ...interface{}) {
	l.Warn(format, v...)
}

func Notice(format string, v ...interface{}) {
	l.Notice(format, v...)
}

func Info(format string, v ...interface{}) {
	l.Info(format, v...)
}

func Debug(format string, v ...interface{}) {
	l.Debug(format, v...)
}
