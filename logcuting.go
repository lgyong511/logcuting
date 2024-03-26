package logcuting

import (
	"os"
	"time"
)

// 配置信息
type Config struct {
	Path       string        //日志文件目录
	Name       string        //日志文件名称
	TimeFormat string        //日志文件名称的时间格式
	Time       time.Duration //日志切割时间间隔，最新切割时间间隔分钟
}

type Logcuting struct {
	Config
	file *os.File  //文件句柄
	t    time.Time //上次日志切割的时间
}

// 创建Logcuting实例
func NewLogcuting(config *Config) *Logcuting {
	l := new(Logcuting)
	l.Config = *config
	// l.Config.Path = config.Name
	// l.Config.Name=config.Name
	// l.Config.TimeFormat=config.TimeFormat
	// l.Config.Time=config.Time
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
	l.cutingByTime()
	return l.file.Write(p)
}

// 根据时间间隔切割日志文件
// 时间为0时，每天0:0:0秒切割
// 有日志输出的时候才会切割
func (l *Logcuting) cutingByTime() {
	l.setTime()

	if l.Time == 0 { // 一天切割一次
		// 获取时间是否到了0时0分0秒
		// 重新创建文件并赋值给l.file

		// 在现在时间上加一天
		t := time.Now().AddDate(0, 0, 1)
		// 判断现在时间的时间戳是否大于等于加一天后的0时0分0秒的时间戳
		if time.Now().Unix() >= time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Unix() {
			l.file.Close()
			l.file, _ = os.OpenFile(l.Path+time.Now().Format(l.TimeFormat)+"_"+l.Name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		}

	} else { //按t间隔切割
		// 计算时间是否过去了t间隔
		// 重新创建文件并赋值给l.file

		// 判断上次日志切割的时间是否大于等于日志切割时间间隔
		if time.Since(l.t) >= l.Time {
			l.file.Close()
			l.file, _ = os.OpenFile(l.Path+time.Now().Format(l.TimeFormat)+"_"+l.Name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			l.t = time.Now()
		}

	}
}

// 根据时间间隔切割日志文件
// 时间为0时，每天0:0:0秒切割
// 没有日志输出也会进行切割
// 未启用
func (l *Logcuting) cuting() {
	l.setTime()
	// 一天切割一次
	if l.Time == 0 {
		// 获取时间是否到了当天的0时0分0秒
		// 重新创建文件并赋值给l.file
		for {
			t := time.Now().AddDate(0, 0, 1)
			if time.Now().Unix() >= time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Unix() {
				l.file.Close()
				l.file, _ = os.OpenFile(l.Path+time.Now().Format(l.TimeFormat)+"_"+l.Name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			}
			time.Sleep(time.Second * 30)
		}

	} else { //按t间隔切割
		// 计算时间是否过去了t间隔
		// 重新创建文件并赋值给l.file
		for {
			if time.Since(l.t) >= l.Time {
				l.file.Close()
				l.file, _ = os.OpenFile(l.Path+time.Now().Format(l.TimeFormat)+"_"+l.Name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
				l.t = time.Now()
			}
			time.Sleep(time.Second * 30)
		}

	}
}

// 日志切割时间间隔小于1分钟，返回1分钟
// 日志切割时间间隔等于0分钟，返回0分钟
func (l *Logcuting) setTime() {
	if l.Time == 0 {
		return
	} else if l.Time < time.Minute {
		l.Time = time.Minute
	}
}
