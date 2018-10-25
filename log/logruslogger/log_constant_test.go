package logruslogger

import (
	"testing"
)

func Test_GetLogLevel(t *testing.T) {
	levelName := "debug"
	level := GetLogLevel(levelName)
	t.Logf("get level %v, real level: %v", levelName, level)
}
