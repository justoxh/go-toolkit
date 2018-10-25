package log

// Logger log interface
type Logger interface {
	Debug(f interface{}, args ...interface{})
	Info(f interface{}, args ...interface{})
	Warn(f interface{}, args ...interface{})
	Printf(f interface{}, args ...interface{})
	Panic(f interface{}, args ...interface{})
	Fatal(f interface{}, args ...interface{})
	Error(f interface{}, args ...interface{})

	Debugln(v ...interface{})
	Infoln(args ...interface{})
	Warnln(args ...interface{})
	Printfln(args ...interface{})
	Panicln(args ...interface{})
	Fatalln(args ...interface{})
	Errorln(args ...interface{})
}
