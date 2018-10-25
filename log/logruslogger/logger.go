package logruslogger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

const (
	defaultLevel       = "debug"
	defaultLogFileName = "all.logs"
)

var (
	logMap      map[string]*LogrusLogger
	getLogMutex sync.Mutex
)

func defaultOptions() *Options {
	return &Options{
		Level:          defaultLevel,
		Depth:          8,
		WithCallerHook: true,
		Formatter:      "text",
		DisableConsole: false,
		Write:          false,
		Path:           os.TempDir(),
		FileName:       defaultLogFileName,
		MaxAge:         time.Duration(24) * time.Hour,
		RotationTime:   time.Duration(7*24) * time.Hour,
		Debug:          false,
	}
}

// Options logger options config
type Options struct {
	Level          string
	Depth          int
	WithCallerHook bool
	Formatter      string // only support json and text

	DisableConsole bool
	Write          bool
	Path           string
	FileName       string

	MaxAge        time.Duration
	RotationCount uint // can`t be set with maxAge at the same time
	RotationTime  time.Duration

	Debug bool // if set true, separate
}

// GetLoggerWithOptions with options config
func GetLoggerWithOptions(logName string, options *Options) *LogrusLogger {
	getLogMutex.Lock()
	defer getLogMutex.Unlock()

	if options == nil {
		options = defaultOptions()
	}

	if logMap == nil {
		logMap = make(map[string]*LogrusLogger)
	}
	curLog, ok := logMap[logName]

	if ok {
		return curLog
	}

	log := logrus.New()

	// get logLevel
	level := options.Level
	if level == "" {
		level = defaultLevel
	}
	logLevel := GetLogLevel(level)
	logDir := options.Path
	if logDir == "" {
		logDir = os.TempDir()
	}

	logFileName := options.FileName
	if logFileName == "" {
		logFileName = defaultLogFileName
	}

	printLog := !options.DisableConsole
	depth := options.Depth
	maxAge := options.MaxAge
	rotationCount := options.RotationCount
	rotationTime := options.RotationTime
	withCallerHook := options.WithCallerHook

	log.SetLevel(logLevel)

	if options.Write {
		storeLogDir := logDir

		err := os.MkdirAll(storeLogDir, os.ModePerm)
		if err != nil {
			panic(fmt.Sprintf("creating log file failed: %s", err.Error()))
		}

		path := filepath.Join(storeLogDir, logFileName)
		writer, err := rotatelogs.New(
			path+".%Y%m%d%H%M",
			rotatelogs.WithClock(rotatelogs.Local),
			rotatelogs.WithMaxAge(time.Duration(maxAge)*time.Hour),
			rotatelogs.WithRotationCount(rotationCount),
			rotatelogs.WithRotationTime(time.Duration(rotationTime)*time.Hour),
		)
		if err != nil {
			panic(fmt.Sprintf("rotatelogs log failed: %s", err.Error()))
		}

		var formatter logrus.Formatter

		formatter = &logrus.TextFormatter{}
		if options.Formatter == "json" {
			formatter = &logrus.JSONFormatter{}
		}
		if options.Debug {
			log.AddHook(lfshook.NewHook(
				lfshook.WriterMap{
					logrus.DebugLevel: writer,
					logrus.InfoLevel:  writer,
					logrus.WarnLevel:  writer,
					logrus.ErrorLevel: writer,
					logrus.FatalLevel: writer,
				},
				formatter,
			))

			defaultLogFilePrex := logFileName + "."
			pathMap := lfshook.PathMap{
				logrus.DebugLevel: fmt.Sprintf("%s/%sdebug", storeLogDir, defaultLogFilePrex),
				logrus.InfoLevel:  fmt.Sprintf("%s/%sinfo", storeLogDir, defaultLogFilePrex),
				logrus.WarnLevel:  fmt.Sprintf("%s/%swarn", storeLogDir, defaultLogFilePrex),
				logrus.ErrorLevel: fmt.Sprintf("%s/%serror", storeLogDir, defaultLogFilePrex),
				logrus.FatalLevel: fmt.Sprintf("%s/%sfatal", storeLogDir, defaultLogFilePrex),
			}
			log.AddHook(lfshook.NewHook(
				pathMap,
				formatter,
			))
		} else {
			log.Out = writer
		}

	} else {
		if printLog {
			log.Out = os.Stdout
		}
	}

	if withCallerHook {
		log.AddHook(&CallerHook{depth: depth, module: logName}) // add caller hook to print caller's file and line number
	}
	curLog = &LogrusLogger{
		log: log,
	}
	logMap[logName] = curLog
	fmt.Printf("register logger %v, store in %v, current loggers: %v\n", logName, logDir, logMap)
	return curLog
}
