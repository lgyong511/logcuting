# Logcuting
文件日志切割库，返回一个满足io.Writer接口的实例。  

## 功能
- 按时间间隔切割日志文件  
- 按日志文件大小切割日志文件  
- 配置信息更新  

## 日志文件切割优先级
- 文件大小  
- 时间间隔  

## 配置信息说明
- Name 日志文件目录，需要带时间格式。例如：./log/demo-%Y%m%d%H%M%S.log  
- Time 日志切割时间间隔，0每天凌晨切割，最小切割间隙1分钟。例如：time.Minute  
- Size 日志文件切割大小，单位MB。例如：5

## 导入Logcuting
```go
go get -u github.com/lgyong511/logcuting
```

## 使用Logcuting
```go
package main

import (
	"time"

	"github.com/lgyong511/logcuting"
	"github.com/sirupsen/logrus"
)

func main() {
	// 设置日志级别
	logrus.SetLevel(logrus.DebugLevel)
	// 设置日志输出格式
	logrus.SetFormatter(&logrus.JSONFormatter{})
	// 设置输出文件名行号函数名
	logrus.SetReportCaller(true)

	// 创建logcuting实例
	logcut, err := logcuting.NewLogcuting(&logcuting.Config{
		Name: "./log/demo-%Y%m%d%H%M.log",
		Time: time.Minute,
		// Size: 1,
	})
	if err != nil {
		panic(err)
	}
	// 将logcut作为logrus的输出目标
	logrus.SetOutput(logcut)
	logrus.Info("将日志输出到logcut")
```