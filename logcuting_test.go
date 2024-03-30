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
	t.Log("new success, NewLogcuting=", l)
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

func TestUpdateConfig(t *testing.T) {
	l := logcuting.NewLogcuting(&logcuting.Config{
		Name: "./logs/demo-%Y%m%d%H%M.log",
		Time: 0,
		Size: 0,
	})
	t.Log("new success, NewLogcuting=", l)
	l.UpdateConfig(&logcuting.Config{
		Name: "./logs/demo-%Y-%m-%d-%H-%M.log",
		Time: 10,
		Size: 10,
	})
	t.Log("UpdateConfig success, UpdateConfig NewLogcuting=", l)
}

func BenchmarkWrite(b *testing.B) {
	l := logcuting.NewLogcuting(&logcuting.Config{
		Name: "./logs/demo-%Y%m%d%H%M.log",
		Time: 0,
		Size: 0,
	})
	s := `{"file":"E:/Myz-Nas/学习笔记/go编程语言/go-base/loging-test/main.go:43","func":"main.logConfig","level":"debug","msg":"使用配置文件创建的Debug","time":""}`
	for i := 0; i < b.N; i++ {
		l.Write([]byte(s))
	}
}
