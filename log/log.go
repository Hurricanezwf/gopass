package log

var (
	l *Logger
)

func Init(conf *LogConf) {
	l = NewLogger().Apply(conf)
}

func Error(format string, v ...interface{}) {
	if l != nil {
		l.Error(format, v...)
	}
}

func Warn(format string, v ...interface{}) {
	if l != nil {
		l.Warn(format, v...)
	}
}

func Notice(format string, v ...interface{}) {
	if l != nil {
		l.Notice(format, v...)
	}
}

func Info(format string, v ...interface{}) {
	if l != nil {
		l.Info(format, v...)
	}
}

func Debug(format string, v ...interface{}) {
	if l != nil {
		l.Debug(format, v...)
	}
}
