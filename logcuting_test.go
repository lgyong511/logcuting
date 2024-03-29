package logcuting_test

import (
	"io"
	"testing"

	"github.com/lgyong511/logcuting"
)

func TestNew(t *testing.T) {
	l := logcuting.NewLogcuting(&logcuting.Config{
		Name: "./logs/demo-%Y%m%d%H%M.log",
		Time: 0,
		Size: 10,
	})
	t.Log("success, NewLogcuting=", l)
}

func TestWriter(t *testing.T) {
	var w io.Writer = logcuting.NewLogcuting(&logcuting.Config{
		Name: "./logs/demo-%Y%m%d%H%M.log",
		Time: 0,
		Size: 10,
	})
	_ = w
}

func TestCloser(t *testing.T) {
	var c io.Closer = logcuting.NewLogcuting(&logcuting.Config{
		Name: "./logs/demo-%Y%m%d%H%M.log",
		Time: 0,
		Size: 10,
	})
	_ = c
}
