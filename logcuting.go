package logcuting

import (
	"os"
	"strings"
	"time"
)

// 配置信息
type Config struct {
	Name string        //日志输出目录字符串，例如："./log/demo-%Y%m%d%H%M%S.log"
	Time time.Duration //日志切割时间间隔
	Size int64         //日志文件切割大小，单位MB
}

// 日志切割信息
type Logcuting struct {
	config    *Config   //配置信息，创建Logcuting和调用UpdateConfig时更新
	file      *os.File  //文件实例，每次切割的时候更新
	oldTime   time.Time //上次日志切割的时间，每次切割的时候更新
	oldLayout string    //创建Logcuting和调用UpdateConfig时更新，传统时间格式：%Y-%m-%d %H:%M:%S，从配置信息中截取
	newLayout string    //创建Logcuting和调用UpdateConfig时更新，go语言时间格式："2006-01-02 15:04:05"，根据oldLayout转换
	name      string    //日志输出文件字符串，每次切割的时候更新。"./log/demo-20240327135202.log"

}

// 创建Logcuting实例
func NewLogcuting(config *Config) (*Logcuting, error) {
	l := new(Logcuting)
	l.config = config
	l.setOldLayout()
	l.setNewLayout()
	l.name = l.getName()

	var err error
	l.file, err = os.OpenFile(l.name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	l.oldTime = time.Now()
	return l, nil
}

// 实现io.Writer接口
func (l *Logcuting) Write(p []byte) (n int, err error) {
	if err = l.cuting(); err != nil {
		return
	}

	return l.file.Write(p)
}

// 更新配置信息
func (l *Logcuting) UpdateConfig(config *Config) {
	l.config = config
	l.setOldLayout()
	l.setNewLayout()
}

// 日志切割
func (l *Logcuting) cuting() (err error) {
	//如果config.Size大于0，就按日志文件大小切割，否则就按时间切割
	if l.config.Size > 0 {
		err = l.cutingBySize()
		if err != nil {
			return
		}
	} else {
		err = l.cutingByTime()
		if err != nil {
			return
		}
	}
	return
}

// 根据时间间隔切割日志文件
// config.Time为0时，每天0点0分切割
// 有日志输出的时候才会切割
func (l *Logcuting) cutingByTime() (err error) {
	l.setTime()

	if l.config.Time == 0 { // 一天切割一次
		// 获取时间是否到了0时0分
		// 重新创建文件并赋值给l.file

		// 在现在时间上加一天
		t := time.Now().AddDate(0, 0, 1)
		// 判断现在时间的时间戳是否大于等于加一天后的0时0分0秒的时间戳
		if time.Now().Unix() >= time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Unix() {
			l.file.Close()
			l.name = l.getName()
			l.file, err = os.OpenFile(l.name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				return err
			}
		}

	} else { //按config.Time间隔切割
		// 计算时间是否过去了config.Time间隔
		// 重新创建文件并赋值给l.file

		// 判断上次日志切割的时间是否大于等于日志切割时间间隔
		if time.Since(l.oldTime) >= l.config.Time {
			l.file.Close()
			l.name = l.getName()
			l.file, err = os.OpenFile(l.name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				return err
			}
			l.oldTime = time.Now()
		}

	}
	return
}

// 按日志文件大小切割日志
func (l *Logcuting) cutingBySize() (err error) {
	size := l.getSize()
	if size >= l.config.Size {
		l.file.Close()
		l.name = l.getName()
		l.file, err = os.OpenFile(l.name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		l.oldTime = time.Now()
	}
	return
}

// 设置config.Time
// 日志切割时间间隔小于1分钟，返回1分钟
// 日志切割时间间隔等于0分钟，返回0分钟
func (l *Logcuting) setTime() {
	if l.config.Time == 0 {
		return
	} else if l.config.Time < time.Minute {
		l.config.Time = time.Minute
	}
}

// 设置oldLayout，传统时间格式
func (l *Logcuting) setOldLayout() {
	i := strings.IndexByte(l.config.Name, '%')
	li := strings.LastIndexByte(l.config.Name, '%') + 1
	l.oldLayout = l.config.Name[i : li+1]
}

// 设置newLayout，go语言的时间格式
func (l *Logcuting) setNewLayout() {
	layout := strings.Replace(l.oldLayout, "%Y", "2006", -1)
	layout = strings.Replace(layout, "%m", "01", -1)
	layout = strings.Replace(layout, "%d", "02", -1)
	layout = strings.Replace(layout, "%H", "15", -1)
	layout = strings.Replace(layout, "%M", "04", -1)
	layout = strings.Replace(layout, "%S", "05", -1)
	l.newLayout = layout
}

// 获取日志输出文件路径
func (l *Logcuting) getName() string {

	time := time.Now().Format(l.newLayout)
	return strings.Replace(l.config.Name, l.oldLayout, time, -1)
}

// 获取日志文件的大小
func (l *Logcuting) getSize() int64 {
	fileInfo, _ := l.file.Stat()
	return fileInfo.Size() / 1024 / 1024
}
