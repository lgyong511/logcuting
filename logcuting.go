package logcuting

import (
	"os"
	"time"
)

// 配置信息
type Config struct {
	Path       string //日志文件目录
	Name       string //日志文件名称
	TimeFormat string //日志文件名称的时间格式
}

type Logcuting struct {
	Config
	file *os.File
	t    time.Time
}

// 创建Logcuting实例
func NewLogcuting(config *Config) *Logcuting {
	l := new(Logcuting)
	var err error
	l.file, err = os.OpenFile(l.Path+time.Now().Format(l.TimeFormat)+"_"+l.Name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	l.t = time.Now()
	return l
}

// 实现io.Writer接口
func (l *Logcuting) Write(p []byte) (n int, err error) {
	return l.file.Write(p)
}

func (l *Logcuting) Cuting(t time.Duration) {
	// 一天切割一次
	if t == 0 {
		// 获取时间是否到了当天的0时0分0秒
		// 重新创建文件并赋值给l.file
		for {
			t := time.Now().AddDate(0, 0, 1)
			if time.Now().Unix() == time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Unix() {
				l.file, _ = os.OpenFile(l.Path+time.Now().Format(l.TimeFormat)+"_"+l.Name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			}
			time.Sleep(time.Second * 30)
		}

	} else { //按t间隔切割
		// 计算时间是否过去了t间隔
		// 重新创建文件并赋值给l.file
		for {
			if time.Since(l.t) >= t {
				l.file, _ = os.OpenFile(l.Path+time.Now().Format(l.TimeFormat)+"_"+l.Name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
				l.t = time.Now()
			}
			time.Sleep(time.Second * 10)
		}

	}
}
