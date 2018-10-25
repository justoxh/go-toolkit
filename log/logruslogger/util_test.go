package logruslogger

import (
	"testing"
)

func Test_FormatLog(t *testing.T) {
	str := FormatLog("uu")
	t.Log(str)

	str = FormatLog("u %v", 123)
	t.Log(str)

	str = FormatLog("u", 123)
	t.Log(str)

	str = FormatLog(123)
	t.Log(str)

	str = FormatLog(123, 456)
	t.Log(str)
}
