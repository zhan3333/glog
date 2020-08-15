package glog

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/zhan3333/glog/hook"
	"os"
	"path/filepath"
)

// logrus 的包装工具
// DefaultFormat 配置默认的日志格式
// LogConfigs 日志通道配置

const (
	PanicLevel logrus.Level = iota
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel
)

const (
	// 单文件驱动
	SINGLE = iota
	// 日驱动
	DAILY
)

type Log struct {
	Driver       uint
	Path         string
	Level        logrus.Level
	Days         int
	LogFormatter logrus.Formatter
	ReportCall   bool
	Hooks        []logrus.Hook
}

var LogConfigs = map[string]Log{}
var DefLogChannel = "default"

// 默认日志格式
var DefaultFormat logrus.Formatter = LocalFormatter{&logrus.JSONFormatter{
	PrettyPrint:       false,
	DisableHTMLEscape: true,
}}
var channels = map[string]*logrus.Logger{}
var openFiles []*os.File

type LocalFormatter struct {
	logrus.Formatter
}

func (u LocalFormatter) Format(e *logrus.Entry) ([]byte, error) {
	e.Time = e.Time.Local()
	return u.Formatter.Format(e)
}

// 获取默认通道
func Def() *logrus.Logger {
	return Channel("default")
}

// 获取指定通道
func Channel(name string) *logrus.Logger {
	if l, ok := channels[name]; ok {
		return l
	}
	if c, ok := LogConfigs[name]; ok {
		channels[name] = configLog(c)
		return channels[name]
	}
	return channels["default"]
}

// 加载所有通道
// 将会更改 logrus.StandardLogger() 的行为
func LoadChannels() {
	// default
	channels["default"] = configDefaultLog()
	// channels
	for name, logConf := range LogConfigs {
		channels[name] = configLog(logConf)
	}
}

// 重载所有通道, 通常用于修改了默认通道或者配置后调用
func ReloadChannels() {
	channels = map[string]*logrus.Logger{}
	for _, f := range openFiles {
		_ = f.Close()
	}
	LoadChannels()
}

func configDefaultLog() *logrus.Logger {
	if logC, ok := LogConfigs[DefLogChannel]; ok {
		l := logrus.StandardLogger()
		config(l, logC)
		return l
	}
	return nil
}

func configLog(logConf Log) *logrus.Logger {
	l := logrus.New()
	config(l, logConf)
	return l
}

func config(l *logrus.Logger, c Log) {
	var err error
	var format logrus.Formatter
	if c.LogFormatter != nil {
		format = c.LogFormatter
	} else {
		format = DefaultFormat
	}
	l.SetLevel(c.Level)
	if c.ReportCall {
		l.SetReportCaller(true)
	}
	// add hooks
	for _, h := range c.Hooks {
		l.AddHook(h)
	}
	if c.Driver == DAILY {
		// 日驱动
		l.AddHook(hook.NewLfsHook(c.Path, 7, format))
	} else {
		l.SetFormatter(format)
		// create dir
		err = os.MkdirAll(filepath.Dir(c.Path), os.ModePerm)
		if err != nil {
			panic(fmt.Sprintf("Create dir %s failed: [%+v]", filepath.Dir(c.Path), err))
		}
		// create log file
		f, err := os.OpenFile(c.Path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(fmt.Sprintf("Create file %s failed: [%+v]", c.Path, err))
		}
		openFiles = append(openFiles, f)
		l.SetOutput(f)
	}
}

// 关闭文件
func Close() {
	for _, f := range openFiles {
		_ = f.Close()
	}
}

func SetLogFields(f *map[string]interface{}) logrus.Fields {
	fields := logrus.Fields{}
	if f != nil {
		for k, v := range *f {
			fields[k] = v
		}
	}
	return fields
}
