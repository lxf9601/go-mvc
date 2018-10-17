package logc

import (
	"testing"
)

func TestLog(t *testing.T) {
	Infof("日志有效%v%v", 1, 2)
}
