package logruslogger

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// LogrusLogger logrus logger
type LogrusLogger struct {
	log *logrus.Logger
}

// GetLogger convert LogrusLogger to *logrus.Logger
func (l *LogrusLogger) GetLogger() *logrus.Logger {
	return l.log
}

// Debug wrapper Debug logger
func (l *LogrusLogger) Debug(f interface{}, args ...interface{}) {
	l.log.Debug(FormatLog(f, args...))
}

// Info wrapper Info logger
func (l *LogrusLogger) Info(f interface{}, args ...interface{}) {
	l.log.Info(FormatLog(f, args...))
}

// Warn wrapper Warn logger
func (l *LogrusLogger) Warn(f interface{}, args ...interface{}) {
	l.log.Warn(FormatLog(f, args...))
}

// Printf wrapper Printf logger
func (l *LogrusLogger) Printf(f interface{}, args ...interface{}) {
	l.log.Print(FormatLog(f, args...))
}

// Panic wrapper Panic logger
func (l *LogrusLogger) Panic(f interface{}, args ...interface{}) {
	l.log.Panic(FormatLog(f, args...))
}

// Fatal wrapper Fatal logger
func (l *LogrusLogger) Fatal(f interface{}, args ...interface{}) {
	l.log.Fatal(FormatLog(f, args...))
}

// Error wrapper Error logger
func (l *LogrusLogger) Error(f interface{}, args ...interface{}) {
	l.log.Error(FormatLog(f, args...))
}

// Debugln wrapper Debugln logger
func (l *LogrusLogger) Debugln(v ...interface{}) {
	l.log.Debug(fmt.Sprintln(v...))
}

// Infoln wrapper Infoln logger
func (l *LogrusLogger) Infoln(args ...interface{}) {
	l.log.Info(fmt.Sprintln(args...))
}

// Warnln wrapper Warnln logger
func (l *LogrusLogger) Warnln(args ...interface{}) {
	l.log.Warn(fmt.Sprintln(args...))
}

// Printfln wrapper Printfln logger
func (l *LogrusLogger) Printfln(args ...interface{}) {
	l.log.Print(fmt.Sprintln(args...))
}

// Panicln wrapper Panicln logger
func (l *LogrusLogger) Panicln(args ...interface{}) {
	l.log.Panic(fmt.Sprintln(args...))
}

// Fatalln wrapper Fatalln logger
func (l *LogrusLogger) Fatalln(args ...interface{}) {
	l.log.Fatal(fmt.Sprintln(args...))
}

// Errorln wrapper Errorln logger
func (l *LogrusLogger) Errorln(args ...interface{}) {
	l.log.Error(fmt.Sprintln(args...))
}
