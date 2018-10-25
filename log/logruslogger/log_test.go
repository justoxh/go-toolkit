package logruslogger

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func Test_Log(t *testing.T) {
	options := &Options{}
	options.WithCallerHook = true
	options.Depth = 8
	options.Write = true
	options.RotationTime = time.Duration(2 * time.Second)
	options.Debug = true

	tempTestLog := "testLogDir201808091811"
	dir, _ := os.Getwd()

	storeLogDir := filepath.Join(dir, tempTestLog)
	err := os.MkdirAll(storeLogDir, os.ModePerm)
	if err != nil {
		panic(fmt.Sprintf("creating log file failed: %s", err.Error()))
	}
	options.Path = storeLogDir

	a := GetLoggerWithOptions("a-logrus", options)

	logger := a.GetLogger()
	t.Logf("get logger %v", logger)

	a.Debug("")
	a.Debug(time.Now())
	a.Debug(123, time.Now())
	a.Debug("test %v", time.Now().UnixNano())
	a.Warn("test %v", time.Now().UnixNano())
	a.Info("test %v", time.Now().UnixNano())
	a.Printf("test %v", time.Now().UnixNano())
	a.Printf("test", time.Now().UnixNano())
	a.Error("test %v", time.Now().UnixNano())

	a.Debugln("test", time.Now().UnixNano())
	a.Warnln("test", time.Now().UnixNano())
	a.Infoln("test", time.Now().UnixNano())
	a.Printfln("test", time.Now().UnixNano())
	a.Printfln("test ", time.Now().UnixNano())
	a.Errorln("test", time.Now().UnixNano())

	optionsB := *options
	optionsB.Depth = -1
	optionsB.Formatter = "json"

	b := GetLoggerWithOptions("b-logrus", &optionsB)

	b.Debug("test %v", time.Now().UnixNano())
	b.Warn("test %v", time.Now().UnixNano())
	b.Info("test %v", time.Now().UnixNano())
	b.Printf("test %v", time.Now().UnixNano())
	b.Error("test %v", time.Now().UnixNano())

	b.Debugln("test", time.Now().UnixNano())
	b.Warnln("test", time.Now().UnixNano())
	b.Infoln("test", time.Now().UnixNano())
	b.Printfln("test", time.Now().UnixNano())
	b.Errorln("test", time.Now().UnixNano())

	optionsC := *options
	optionsC.Depth = 19990009900
	c := GetLoggerWithOptions("c-logrus", &optionsC)

	c.Debug("test %v", time.Now().UnixNano())
	c.Warn("test %v", time.Now().UnixNano())

	d := GetLoggerWithOptions("d-logrus", nil)
	d.Debug("123")

	optionsE := *options
	optionsE.DisableConsole = true
	_ = GetLoggerWithOptions("e-logrus", &optionsE)

	optionsF := *options
	optionsF.Path = ""
	_ = GetLoggerWithOptions("f-logrus", &optionsF)

	os.RemoveAll(options.Path)
	t.Logf("clear temp test log files")
}

func Test_PanicLog(t *testing.T) {
	options := &Options{}
	options.WithCallerHook = true
	options.Depth = 8
	options.Write = true

	a := GetLoggerWithOptions("a-logrus", options)
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("fatal log, err: %v\n", err)

			defer func() {
				if err := recover(); err != nil {
					fmt.Printf("fatal log, err: %v\n", err)
				}
			}()
			a.Panicln("fatal test")
		}
	}()

	a.Panic("fatal test")
}

func Test_ErrorPath(t *testing.T) {
	options := &Options{}
	options.WithCallerHook = true
	options.Depth = 8
	options.Write = true
	options.Path = "/logtest"

	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("fatal log, err: %v\n", err)

			defer func() {
				if err := recover(); err != nil {
					fmt.Printf("fatal log, err: %v\n", err)
				}
			}()
			options.Path = ""
			options.RotationCount = 1
			options.MaxAge = 1111
			_ = GetLoggerWithOptions("aa-logrus", options)

		}
	}()
	_ = GetLoggerWithOptions("aa-logrus", options)

}
