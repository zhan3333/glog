package glog_test

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/zhan3333/glog"
	"log"
	"testing"
)

func TestMain(m *testing.M) {
	glog.LogConfigs = map[string]glog.Log{
		glog.DefLogChannel: {
			Driver:       glog.DAILY,
			Path:         "logs/def.log",
			Level:        glog.DebugLevel,
			Days:         30,
			LogFormatter: nil,
			ReportCall:   false,
			Hooks:        nil,
		},
	}
	m.Run()
}

func TestLog(t *testing.T) {
	log.Print("log by log")
	glog.Def().Print("log by glog.Default()")
	glog.Channel("gin").Print("log by glog.Channel(\"gin\")")
	glog.Channel("gin").Print("log by glog.Channel(\"gin\")")
}

func TestAllChannel(t *testing.T) {
	for name := range glog.LogConfigs {
		glog.Channel(name).Printf("Test")
	}
}

type Format struct {
}

func (Format) Format(entry *logrus.Entry) ([]byte, error) {
	return json.Marshal(entry.Message)
}

func TestSetDefaultFormatter(t *testing.T) {
	newFormat := Format{}
	glog.LogConfigs = map[string]glog.Log{
		glog.DefLogChannel: {
			Driver:       glog.DAILY,
			Path:         "logs/def.log",
			Level:        glog.DebugLevel,
			Days:         30,
			LogFormatter: nil,
			ReportCall:   false,
			Hooks:        nil,
		},
	}
	glog.DefaultFormat = newFormat
	glog.ReloadChannels()
	glog.Def().WithFields(logrus.Fields{
		"test": "123",
	}).Info("test")
}
