package glog_test

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/zhan3333/glog"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestLog(t *testing.T) {
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
	log.Print("log by log")
	glog.Def().Print("log by glog.Default()")
	glog.Channel("gin").Print("log by glog.Channel(\"gin\")")
	glog.Channel("gin").Print("log by glog.Channel(\"gin\")")
}

func TestAllChannel(t *testing.T) {
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

func TestOut(t *testing.T) {
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
	glog.LogConfigs = map[string]glog.Log{
		glog.DefLogChannel: {
			Driver:       glog.DAILY,
			Path:         "logs/test.log",
			Level:        glog.DebugLevel,
			Days:         30,
			LogFormatter: nil,
			ReportCall:   false,
			Hooks:        nil,
		},
	}
	glog.ReloadChannels()
	// channel write
	glog.Def().WithFields(logrus.Fields{
		"test": "log from glog.Def().Info()",
	}).Info("test")
	// out obj write
	_, err := glog.Def().Out.Write([]byte("log from glog.Def().Out"))
	assert.Nil(t, err)
}

func TestWrite(t *testing.T) {
	glog.LogConfigs = map[string]glog.Log{
		glog.DefLogChannel: {
			Driver: glog.DAILY,
			Path:   "logs/test.log",
			Level:  glog.DebugLevel,
			Days:   30,
		},
	}
	glog.ReloadChannels()
	glog.Def().WithFields(logrus.Fields{
		"test": "write from glog.Def()",
	}).Info("test")
	glog.Def().Write.Write([]byte("write from glog.Def().Write"))
}

func TestDefaultChannel(t *testing.T) {
	glog.Def().Infoln("test")
}

// 测试单文件驱动
func TestSingleDriver(t *testing.T) {
	assert.Nil(t, os.Remove("logs/test.log"))
	glog.LogConfigs = map[string]glog.Log{
		glog.DefLogChannel: {
			Driver: glog.SINGLE,
			Path:   "logs/test.log",
			Level:  glog.DebugLevel,
			Days:   30,
		},
	}
	glog.Def().Infoln("test")
	b, err := ioutil.ReadFile("logs/test.log")
	assert.Nil(t, err)
	t.Log(string(b))
	type Msg struct {
		Msg string
	}
	var msg Msg
	assert.Nil(t, json.Unmarshal(b, &msg))
	assert.Equal(t, "test", msg.Msg)
}
