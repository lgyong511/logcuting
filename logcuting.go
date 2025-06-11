package logcuting

import (
	"os"
	"path/filepath"
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
	config    *Config  //配置信息，创建Logcuting和调用UpdateConfig时更新
	file      *os.File //文件实例，每次切割的时候更新
	oldTime   int64    //上次日志切割的时间，每次切割的时候更新
	oldLayout string   //创建Logcuting和调用UpdateConfig时更新，传统时间格式：%Y-%m-%d %H:%M:%S，从配置信息中截取
	newLayout string   //创建Logcuting和调用UpdateConfig时更新，go语言时间格式："2006-01-02 15:04:05"，根据oldLayout转换
	name      string   //日志输出文件字符串，每次切割的时候更新。"./log/demo-20240327135202.log"

}

// 创建Logcuting实例
func NewLogcuting(config *Config) *Logcuting {
	l := new(Logcuting)
	l.config = config

	l.setLayouts()

	l.name = l.getName()

	l.file = nil
	l.oldTime = time.Now().UnixMicro()
	return l
}

// 实现io.Writer接口
func (l *Logcuting) Write(p []byte) (n int, err error) {
	if err = l.cuting(); err != nil {
		return
	}

	return l.file.Write(p)
}

// 实现io.Close接口
func (l *Logcuting) Close() error {
	l.config = nil
	l.oldTime = 0
	l.oldLayout = ""
	l.name = ""
	return l.file.Close()
}

// 更新配置信息
func (l *Logcuting) UpdateConfig(config *Config) {
	l.config = config

	l.setLayouts()
}

// 日志切割
func (l *Logcuting) cuting() (err error) {
	// 判断是否打开文件和目录是否存在
	if l.file == nil {
		if err := os.MkdirAll(filepath.Dir(l.name), 0755); err != nil {
			return err
		}
		if err = l.openFile(); err != nil {
			return err
		}
	}
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

			if err := l.rotateFile(); err != nil {
				return err
			}
		}
	} else { //按config.Time间隔切割
		// 计算时间是否过去了config.Time间隔
		// 重新创建文件并赋值给l.file

		// 判断上次日志切割的时间是否大于等于日志切割时间间隔
		if time.Since(time.UnixMicro(l.oldTime)) >= l.config.Time {

			if err := l.rotateFile(); err != nil {
				return err
			}
			l.oldTime = time.Now().UnixMicro()
		}
	}
	return
}

// 按日志文件大小切割日志
func (l *Logcuting) cutingBySize() (err error) {
	size := l.getSize()
	if size >= l.config.Size {

		if err := l.rotateFile(); err != nil {
			return err
		}
		l.oldTime = time.Now().UnixMicro()
	}
	return
}

// rotateFile 执行文件切割和新建操作
func (l *Logcuting) rotateFile() error {
	if err := l.file.Close(); err != nil {
		return err
	}
	l.name = l.getName()
	return l.openFile()
}

// 设置file
// 根据配置信息打开文件并赋值给file
func (l *Logcuting) openFile() (err error) {
	l.file, err = os.OpenFile(l.name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
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

// 设置oldLayout，传统时间格式，从config.Name截取
func (l *Logcuting) setOldLayout() {
	i := strings.IndexByte(l.config.Name, '%')
	li := strings.LastIndexByte(l.config.Name, '%')

	// i和li返回值是-1，表示config.Name中不包含%Y%m%d%H%M%S样式
	if i != -1 { // 只要i和li其中一个的返回值不是-1，i和li都不可能是-1
		l.oldLayout = l.config.Name[i : li+2]
	} else {
		l.oldLayout = "%Y%m%d%H%M"
	}
}

// 设置newLayout，go语言的时间格式，用oldLayout替换
func (l *Logcuting) setNewLayout() {

	replacements := []struct {
		old, new string
	}{
		{"%Y", "2006"},
		{"%m", "01"},
		{"%d", "02"},
		{"%H", "15"},
		{"%M", "04"},
		{"%S", "05"},
	}

	layout := l.oldLayout
	for _, r := range replacements {
		layout = strings.ReplaceAll(layout, r.old, r.new)
	}
	l.newLayout = layout

}

// setLayouts 设置时间格式
func (l *Logcuting) setLayouts() {
	l.setOldLayout()
	l.setNewLayout()
}

// 获取日志输出文件路径
func (l *Logcuting) getName() string {
	t := time.Now().Format(l.newLayout)
	return strings.ReplaceAll(l.config.Name, l.oldLayout, t)
}

// 获取日志文件的大小
func (l *Logcuting) getSize() int64 {
	fileInfo, _ := l.file.Stat()
	return fileInfo.Size() / 1024 / 1024
}
